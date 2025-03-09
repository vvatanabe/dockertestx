package redis_test

import (
	"context"
	redis2 "github.com/vvatanabe/dockertestx/redis"
	"github.com/vvatanabe/dockertestx/sql"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
)

// TestDefaultRedis demonstrates using NewRedis with default options.
func TestDefaultRedis(t *testing.T) {
	// Start a Redis container with default options.
	client, cleanup := redis2.NewRedis(t)
	defer cleanup()

	ctx := context.Background()

	// Test setting and getting a value
	key := "test-key"
	value := "test-value"
	items := map[string]interface{}{
		key: value,
	}

	// Set the item
	if err := redis2.PrepRedis(t, client, items, time.Hour); err != nil {
		t.Fatalf("failed to set item: %v", err)
	}

	// Get the item back
	got, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed to get item: %v", err)
	}

	// Verify the value
	if got != value {
		t.Errorf("expected value '%s', but got '%s'", value, got)
	}
}

// TestRedisWithCustomRunOptions demonstrates overriding default RunOptions.
func TestRedisWithCustomRunOptions(t *testing.T) {
	// Custom RunOption to override the default tag
	customTag := func(opts *dockertest.RunOptions) {
		opts.Tag = "7.0" // Use a specific version
	}

	// Start a Redis container with a custom tag
	client, cleanup := redis2.NewRedisWithOptions(t, []sql.RunOption{customTag})
	defer cleanup()

	ctx := context.Background()

	// Test basic functionality
	key := "custom-test-key"
	value := "custom-test-value"
	items := map[string]interface{}{
		key: value,
	}

	if err := redis2.PrepRedis(t, client, items, time.Hour); err != nil {
		t.Fatalf("failed to set item: %v", err)
	}

	got, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed to get item: %v", err)
	}

	if got != value {
		t.Errorf("expected value '%s', but got '%s'", value, got)
	}
}

// TestRedisWithCustomHostOptions demonstrates providing host configuration options.
func TestRedisWithCustomHostOptions(t *testing.T) {
	// Host option to set AutoRemove to true
	autoRemove := func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	}

	// Start a Redis container with AutoRemove option
	client, cleanup := redis2.NewRedisWithOptions(t, nil, autoRemove)
	defer cleanup()

	ctx := context.Background()

	// Test multiple operations
	items := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	// Use PrepRedis to set up test data
	if err := redis2.PrepRedis(t, client, items, time.Hour); err != nil {
		t.Fatalf("PrepRedis failed: %v", err)
	}

	// Verify all items
	for key, value := range items {
		got, err := client.Get(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to get item '%s': %v", key, err)
		}
		if got != value {
			t.Errorf("expected value '%v' for key '%s', but got '%v'", value, key, got)
		}
	}
}

// TestRedisDataTypes demonstrates various Redis data type operations.
func TestRedisDataTypes(t *testing.T) {
	client, cleanup := redis2.NewRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("String", func(t *testing.T) {
		items := map[string]interface{}{
			"string1": "value1",
			"string2": "value2",
		}
		if err := redis2.PrepRedis(t, client, items, time.Hour); err != nil {
			t.Fatalf("failed to set strings: %v", err)
		}

		for key, value := range items {
			got, err := client.Get(ctx, key).Result()
			if err != nil {
				t.Fatalf("failed to get string '%s': %v", key, err)
			}
			if got != value {
				t.Errorf("expected string '%v', got '%v'", value, got)
			}
		}
	})

	t.Run("List", func(t *testing.T) {
		key := "list1"
		values := []interface{}{"item1", "item2", "item3"}

		if err := redis2.PrepRedisList(t, client, key, values); err != nil {
			t.Fatalf("failed to create list: %v", err)
		}

		got, err := client.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			t.Fatalf("failed to get list: %v", err)
		}

		if len(got) != len(values) {
			t.Errorf("expected list length %d, got %d", len(values), len(got))
		}

		for i, value := range values {
			if got[i] != value {
				t.Errorf("expected list item '%v', got '%v'", value, got[i])
			}
		}
	})

	t.Run("Hash", func(t *testing.T) {
		key := "hash1"
		fields := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}

		if err := redis2.PrepRedisHash(t, client, key, fields); err != nil {
			t.Fatalf("failed to create hash: %v", err)
		}

		got, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to get hash: %v", err)
		}

		if len(got) != len(fields) {
			t.Errorf("expected hash length %d, got %d", len(fields), len(got))
		}

		for field, value := range fields {
			if got[field] != value {
				t.Errorf("expected hash field '%s' value '%v', got '%v'", field, value, got[field])
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		key := "set1"
		members := []interface{}{"member1", "member2", "member3"}

		if err := redis2.PrepRedisSet(t, client, key, members); err != nil {
			t.Fatalf("failed to create set: %v", err)
		}

		got, err := client.SMembers(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to get set members: %v", err)
		}

		if len(got) != len(members) {
			t.Errorf("expected set size %d, got %d", len(members), len(got))
		}

		for _, member := range members {
			exists, err := client.SIsMember(ctx, key, member).Result()
			if err != nil {
				t.Fatalf("failed to check set membership: %v", err)
			}
			if !exists {
				t.Errorf("expected member '%v' to exist in set", member)
			}
		}
	})

	t.Run("SortedSet", func(t *testing.T) {
		key := "zset1"
		members := map[string]float64{
			"member1": 1.0,
			"member2": 2.0,
			"member3": 3.0,
		}

		if err := redis2.PrepRedisSortedSet(t, client, key, members); err != nil {
			t.Fatalf("failed to create sorted set: %v", err)
		}

		got, err := client.ZRangeWithScores(ctx, key, 0, -1).Result()
		if err != nil {
			t.Fatalf("failed to get sorted set: %v", err)
		}

		if len(got) != len(members) {
			t.Errorf("expected sorted set size %d, got %d", len(members), len(got))
		}

		for _, z := range got {
			member := z.Member.(string)
			score := z.Score
			expectedScore := members[member]
			if score != expectedScore {
				t.Errorf("expected score %f for member '%s', got %f", expectedScore, member, score)
			}
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		key := "expiring-key"
		value := "expiring-value"
		items := map[string]interface{}{
			key: value,
		}

		// Set with 1 second expiration
		if err := redis2.PrepRedis(t, client, items, time.Second); err != nil {
			t.Fatalf("failed to set expiring item: %v", err)
		}

		// Verify it exists
		got, err := client.Get(ctx, key).Result()
		if err != nil {
			t.Fatalf("failed to get item: %v", err)
		}
		if got != value {
			t.Errorf("expected value '%s', got '%s'", value, got)
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Verify it's gone
		_, err = client.Get(ctx, key).Result()
		if err != redis.Nil {
			t.Errorf("expected key to be expired, but got: %v", err)
		}
	})
}
