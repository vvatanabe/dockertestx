package dockertestx

import (
	"fmt"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	defaultMemcachedImage = "memcached"
	defaultMemcachedTag   = "1.6.18"
)

// NewMemcached starts a Memcached Docker container using the default settings and returns a connected
// *memcache.Client along with a cleanup function. It uses the default Memcached image ("memcached")
// with tag "1.6.18". For more customization, use NewMemcachedWithOptions.
func NewMemcached(t testing.TB) (*memcache.Client, func()) {
	return NewMemcachedWithOptions(t, nil)
}

// NewMemcachedWithOptions starts a Memcached Docker container using Docker and returns a connected
// *memcache.Client along with a cleanup function. It applies the default settings:
//   - Repository: "memcached"
//   - Tag: "1.6.18"
//
// Additional RunOption functions can be provided via the runOpts parameter to override these defaults,
// and optional host configuration functions can be provided via hostOpts.
func NewMemcachedWithOptions(t testing.TB, runOpts []RunOption, hostOpts ...func(*docker.HostConfig)) (*memcache.Client, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to connect to docker: %s", err)
	}

	// Set default run options for Memcached
	defaultRunOpts := &dockertest.RunOptions{
		Repository: defaultMemcachedImage,
		Tag:        defaultMemcachedTag,
	}

	// Apply any provided RunOption functions to override defaults
	for _, opt := range runOpts {
		opt(defaultRunOpts)
	}

	// Pass optional host configuration options
	resource, err := pool.RunWithOptions(defaultRunOpts, hostOpts...)
	if err != nil {
		t.Fatalf("failed to start memcached container: %s", err)
	}

	actualPort := resource.GetHostPort("11211/tcp")
	if actualPort == "" {
		_ = pool.Purge(resource)
		t.Fatal("no host port was assigned for the memcached container")
	}
	t.Logf("memcached container is running on host port '%s'", actualPort)

	// Create Memcached client
	var client *memcache.Client

	// Try to connect to Memcached with retries
	if err = pool.Retry(func() error {
		client = memcache.New(actualPort)
		// Ping the server by attempting to get a non-existent key
		// This will return ErrCacheMiss if the server is responsive
		_, err := client.Get("test-connection")
		if err != nil && err != memcache.ErrCacheMiss {
			return fmt.Errorf("failed to connect to memcached: %w", err)
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("could not connect to memcached: %s", err)
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Logf("failed to remove memcached container: %s", err)
		}
	}

	return client, cleanup
}

// PrepMemcached sets up test data in a Memcached instance.
// It accepts a list of memcache.Item pointers and stores them in the cache.
// If any operation fails, it returns an error.
func PrepMemcached(t testing.TB, client *memcache.Client, items ...*memcache.Item) error {
	t.Helper()

	for _, item := range items {
		if item.Expiration == 0 {
			// Set default expiration to 1 hour if not specified
			item.Expiration = int32(time.Hour.Seconds())
		}
		if err := client.Set(item); err != nil {
			return fmt.Errorf("failed to set item with key '%s': %w", item.Key, err)
		}
	}
	return nil
}
