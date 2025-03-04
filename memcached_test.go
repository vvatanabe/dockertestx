package dockertestx_test

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/vvatanabe/dockertestx"
)

// TestDefaultMemcached demonstrates using NewMemcached with default options.
func TestDefaultMemcached(t *testing.T) {
	// Start a Memcached container with default options.
	client, cleanup := dockertestx.NewMemcached(t)
	defer cleanup()

	// Test setting and getting a value
	key := "test-key"
	value := []byte("test-value")
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}

	// Set the item
	if err := client.Set(item); err != nil {
		t.Fatalf("failed to set item: %v", err)
	}

	// Get the item back
	got, err := client.Get(key)
	if err != nil {
		t.Fatalf("failed to get item: %v", err)
	}

	// Verify the value
	if string(got.Value) != string(value) {
		t.Errorf("expected value '%s', but got '%s'", value, got.Value)
	}
}

// TestMemcachedWithCustomRunOptions demonstrates overriding default RunOptions.
func TestMemcachedWithCustomRunOptions(t *testing.T) {
	// Custom RunOption to override the default tag
	customTag := func(opts *dockertest.RunOptions) {
		opts.Tag = "1.6.17" // Use a specific version
	}

	// Start a Memcached container with a custom tag
	client, cleanup := dockertestx.NewMemcachedWithOptions(t, []dockertestx.RunOption{customTag})
	defer cleanup()

	// Test basic functionality
	key := "custom-test-key"
	value := []byte("custom-test-value")
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}

	if err := client.Set(item); err != nil {
		t.Fatalf("failed to set item: %v", err)
	}

	got, err := client.Get(key)
	if err != nil {
		t.Fatalf("failed to get item: %v", err)
	}

	if string(got.Value) != string(value) {
		t.Errorf("expected value '%s', but got '%s'", value, got.Value)
	}
}

// TestMemcachedWithCustomHostOptions demonstrates providing host configuration options.
func TestMemcachedWithCustomHostOptions(t *testing.T) {
	// Host option to set AutoRemove to true
	autoRemove := func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	}

	// Start a Memcached container with AutoRemove option
	client, cleanup := dockertestx.NewMemcachedWithOptions(t, nil, autoRemove)
	defer cleanup()

	// Test multiple operations
	items := []*memcache.Item{
		{
			Key:   "key1",
			Value: []byte("value1"),
		},
		{
			Key:   "key2",
			Value: []byte("value2"),
		},
	}

	// Use PrepMemcached to set up test data
	if err := dockertestx.PrepMemcached(t, client, items...); err != nil {
		t.Fatalf("PrepMemcached failed: %v", err)
	}

	// Verify all items
	for _, item := range items {
		got, err := client.Get(item.Key)
		if err != nil {
			t.Fatalf("failed to get item '%s': %v", item.Key, err)
		}
		if string(got.Value) != string(item.Value) {
			t.Errorf("expected value '%s' for key '%s', but got '%s'",
				item.Value, item.Key, got.Value)
		}
	}
}

// TestMemcachedOperations demonstrates various Memcached operations.
func TestMemcachedOperations(t *testing.T) {
	client, cleanup := dockertestx.NewMemcached(t)
	defer cleanup()

	t.Run("Set and Get", func(t *testing.T) {
		item := &memcache.Item{
			Key:   "test1",
			Value: []byte("value1"),
		}
		if err := client.Set(item); err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		if got, err := client.Get("test1"); err != nil {
			t.Fatalf("failed to get: %v", err)
		} else if string(got.Value) != "value1" {
			t.Errorf("expected 'value1', got '%s'", got.Value)
		}
	})

	t.Run("Add", func(t *testing.T) {
		item := &memcache.Item{
			Key:   "test2",
			Value: []byte("value2"),
		}
		if err := client.Add(item); err != nil {
			t.Fatalf("failed to add: %v", err)
		}
		// Adding same key should fail
		if err := client.Add(item); err != memcache.ErrNotStored {
			t.Errorf("expected ErrNotStored, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test3"
		item := &memcache.Item{
			Key:   key,
			Value: []byte("value3"),
		}
		if err := client.Set(item); err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		if err := client.Delete(key); err != nil {
			t.Fatalf("failed to delete: %v", err)
		}
		if _, err := client.Get(key); err != memcache.ErrCacheMiss {
			t.Errorf("expected ErrCacheMiss, got %v", err)
		}
	})

	t.Run("Increment", func(t *testing.T) {
		key := "counter"
		item := &memcache.Item{
			Key:   key,
			Value: []byte("1"),
		}
		if err := client.Set(item); err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		if newVal, err := client.Increment(key, 1); err != nil {
			t.Fatalf("failed to increment: %v", err)
		} else if newVal != 2 {
			t.Errorf("expected 2, got %d", newVal)
		}
	})
}
