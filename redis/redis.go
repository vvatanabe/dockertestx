package redis

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"github.com/vvatanabe/dockertestx"
	"testing"
	"time"
)

const (
	defaultRedisImage = "redis"
	defaultRedisTag   = "7.2"
)

// NewRedis starts a Redis Docker container using the default settings and returns a connected
// *redis.Client along with a cleanup function. It uses the default Redis image ("redis")
// with tag "7.2". For more customization, use NewRedisWithOptions.
func NewRedis(t testing.TB) (*redis.Client, func()) {
	return NewRedisWithOptions(t, nil)
}

// NewRedisWithOptions starts a Redis Docker container using Docker and returns a connected
// *redis.Client along with a cleanup function. It applies the default settings:
//   - Repository: "redis"
//   - Tag: "7.2"
//
// Additional RunOption functions can be provided via the runOpts parameter to override these defaults,
// and optional host configuration functions can be provided via hostOpts.
func NewRedisWithOptions(t testing.TB, runOpts []dockertestx.RunOption, hostOpts ...func(*docker.HostConfig)) (*redis.Client, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to connect to docker: %s", err)
	}

	// Set default run options for Redis
	defaultRunOpts := &dockertest.RunOptions{
		Repository: defaultRedisImage,
		Tag:        defaultRedisTag,
	}

	// Apply any provided RunOption functions to override defaults
	for _, opt := range runOpts {
		opt(defaultRunOpts)
	}

	// Pass optional host configuration options
	resource, err := pool.RunWithOptions(defaultRunOpts, hostOpts...)
	if err != nil {
		t.Fatalf("failed to start redis container: %s", err)
	}

	actualPort := resource.GetHostPort("6379/tcp")
	if actualPort == "" {
		_ = pool.Purge(resource)
		t.Fatal("no host port was assigned for the redis container")
	}
	t.Logf("redis container is running on host port '%s'", actualPort)

	// Create Redis client
	var client *redis.Client

	// Try to connect to Redis with retries
	ctx := context.Background()
	if err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: actualPort,
		})
		// Ping the server to check if it's responsive
		return client.Ping(ctx).Err()
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("could not connect to redis: %s", err)
	}

	cleanup := func() {
		if err := client.Close(); err != nil {
			t.Logf("failed to close Redis client: %s", err)
		}
		if err := pool.Purge(resource); err != nil {
			t.Logf("failed to remove redis container: %s", err)
		}
	}

	return client, cleanup
}

// PrepRedis sets up test data in a Redis instance.
// It accepts a map of key-value pairs and stores them in the cache.
// If any operation fails, it returns an error.
func PrepRedis(t testing.TB, client *redis.Client, items map[string]interface{}, expiration time.Duration) error {
	t.Helper()

	ctx := context.Background()
	for key, value := range items {
		if err := client.Set(ctx, key, value, expiration).Err(); err != nil {
			return fmt.Errorf("failed to set item with key '%s': %w", key, err)
		}
	}
	return nil
}

// PrepRedisList sets up test data in a Redis list.
// It accepts a key and a list of values to be stored.
// If any operation fails, it returns an error.
func PrepRedisList(t testing.TB, client *redis.Client, key string, values []interface{}) error {
	t.Helper()

	ctx := context.Background()
	if err := client.RPush(ctx, key, values...).Err(); err != nil {
		return fmt.Errorf("failed to push items to list '%s': %w", key, err)
	}
	return nil
}

// PrepRedisHash sets up test data in a Redis hash.
// It accepts a key and a map of field-value pairs to be stored.
// If any operation fails, it returns an error.
func PrepRedisHash(t testing.TB, client *redis.Client, key string, fields map[string]interface{}) error {
	t.Helper()

	ctx := context.Background()
	if err := client.HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("failed to set hash fields for key '%s': %w", key, err)
	}
	return nil
}

// PrepRedisSet sets up test data in a Redis set.
// It accepts a key and a list of members to be stored.
// If any operation fails, it returns an error.
func PrepRedisSet(t testing.TB, client *redis.Client, key string, members []interface{}) error {
	t.Helper()

	ctx := context.Background()
	if err := client.SAdd(ctx, key, members...).Err(); err != nil {
		return fmt.Errorf("failed to add members to set '%s': %w", key, err)
	}
	return nil
}

// PrepRedisSortedSet sets up test data in a Redis sorted set.
// It accepts a key and a map of member-score pairs to be stored.
// If any operation fails, it returns an error.
func PrepRedisSortedSet(t testing.TB, client *redis.Client, key string, members map[string]float64) error {
	t.Helper()

	ctx := context.Background()
	for member, score := range members {
		if err := client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err(); err != nil {
			return fmt.Errorf("failed to add member '%s' to sorted set '%s': %w", member, key, err)
		}
	}
	return nil
}
