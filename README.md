# dockertestx

**dockertestx** is a Go testing library that leverages [dockertest](https://github.com/ory/dockertest) to simplify integration testing with databases and other services using Docker containers. Originally forked from [sqltest](https://github.com/vvatanabe/sqltest), it extends the functionality to support not only SQL databases (MySQL, PostgreSQL) but also other data stores and services, making it easier to set up and manage Docker containers for your integration tests.

## Features

### Supported Services
- **SQL Databases**: MySQL and PostgreSQL container management
- **Cache Services**: Memcached and Redis support
- **Object Storage**: MinIO (S3-compatible) support
- **Future Support**: MongoDB and other data stores
- **Extensibility**: Easy to add custom service containers

### Simple & Powerful
- **Easy to Use**: Clean and intuitive API
- **Automatic**: Container lifecycle management and health checks
- **Flexible**: Rich customization options and configurations
- **Reliable**: Built-in test helpers and utilities

## Installation

Use `go get` to install the package:

```bash
go get github.com/vvatanabe/dockertestx
```

## Usage

The package provides several key functions:

### NewMySQL
This function starts a MySQL Docker container using default settings. It uses the MySQL image (`"mysql"`) with the default tag (`"8.0"`). It returns a connected `*sql.DB` instance along with a cleanup function that ensures the container is removed after the test completes. For most cases, you can use NewMySQL directly for a quick setup.

### NewMySQLWithOptions
For advanced usage, `NewMySQLWithOptions` allows you to customize the container's settings. In addition to the defaults used by `NewMySQL`, you can pass one or more `RunOption` functions to override any default configuration (for example, changing the environment variables, command, mounts, etc.).
You can also provide optional host configuration options (via variadic functions) that allow you to adjust Docker's `HostConfig` settings (e.g., setting `AutoRemove` to true).

### NewPostgres
This function starts a PostgreSQL Docker container using default settings. It uses the PostgreSQL image (`"postgres"`) with the default tag (`"13"`). It returns a connected *sql.DB and a cleanup function that removes the container after the test is done. For most cases, you can use NewPostgres directly for a quick setup.

### NewPostgresWithOptions
Similar to the MySQL variant, `NewPostgresWithOptions` allows you to override the default settings by accepting additional `RunOption` functions. You can customize the container configuration (e.g., changing environment variables or other run options) and supply optional host configuration functions to adjust Docker's `HostConfig` (such as setting `AutoRemove`).

### NewMemcached
This function starts a Memcached Docker container using default settings. It uses the Memcached image (`"memcached"`) with the default tag (`"1.6.18"`). It returns a connected `*memcache.Client` along with a cleanup function that removes the container after the test is done.

### NewMemcachedWithOptions
Similar to other services, `NewMemcachedWithOptions` allows you to customize the container's settings through `RunOption` functions and host configuration options. This provides flexibility in configuring the Memcached container for specific test scenarios.

### PrepMemcached
A helper function that sets up test data in a Memcached instance. It accepts a list of `memcache.Item` pointers and stores them in the cache with optional expiration times.

### NewRedis
This function starts a Redis Docker container using default settings. It uses the Redis image (`"redis"`) with the default tag (`"7.2"`). It returns a connected `*redis.Client` along with a cleanup function that removes the container after the test is done.

### NewRedisWithOptions
Similar to other services, `NewRedisWithOptions` allows you to customize the container's settings through `RunOption` functions and host configuration options. This provides flexibility in configuring the Redis container for specific test scenarios.

### PrepRedis
A helper function that sets up test data in a Redis instance. It accepts a map of key-value pairs and stores them in the cache with optional expiration times.

### PrepRedisList
A helper function that sets up list data in a Redis instance. It accepts a key and a list of values to be stored.

### PrepRedisHash
A helper function that sets up hash data in a Redis instance. It accepts a key and a map of field-value pairs to be stored.

### PrepRedisSet
A helper function that sets up set data in a Redis instance. It accepts a key and a list of members to be stored.

### PrepRedisSortedSet
A helper function that sets up sorted set data in a Redis instance. It accepts a key and a map of member-score pairs to be stored.

### NewDockerDB
A helper function that starts a Docker container with the given run options, waits for the database to be ready, and returns a connected `*sql.DB` along with a cleanup function.

### PrepDatabase
Prepares the test database by executing provided schema (DDL) and initial data (DML) SQL statements. The initial data insertion is performed within a transaction to ensure consistency.

### InitialDBSetup
A helper struct used with `PrepDatabase` to specify the schema and initial data for setting up your test database.

### NewMinIO
This function starts a MinIO Docker container using default settings. It uses the MinIO image (`"minio/minio"`) with the default tag (`"latest"`). It returns a configured `*s3.Client` from AWS SDK Go v2 along with a cleanup function that ensures the container is removed after the test completes. For most cases, you can use NewMinIO directly for a quick setup.

### NewMinIOWithOptions
Similar to other services, `NewMinIOWithOptions` allows you to customize the container's settings through `RunOption` functions and host configuration options. It applies the default settings:
  - Repository: "minio/minio"
  - Tag: "latest"
  - Environment: MINIO_ROOT_USER=minioadmin, MINIO_ROOT_PASSWORD=minioadmin
  - Command: ["server", "/data"]

This provides flexibility in configuring the MinIO container for specific test scenarios.

### PrepBucket
A helper function that creates a bucket in MinIO if it doesn't already exist. This is useful for setting up the necessary storage structure before uploading objects.

### UploadObject
A helper function that uploads a single object (file) to a specified bucket in MinIO. It accepts a bucket name, key (object path), and the object data as a byte slice.

### PrepS3Objects
A helper function that prepares a bucket with multiple objects. It first ensures the bucket exists (creating it if necessary), then uploads all the specified objects. It accepts a bucket name and a map of object keys to byte slices containing the object data.

## Examples

### MySQL Example
[Previous MySQL example code remains the same]

### PostgreSQL Example
[Previous PostgreSQL example code remains the same]

### Memcached Example
[Previous Memcached example code remains the same]

### Redis Example

```go
package dockertestx_test

import (
    "context"
    "testing"
    "time"
    "github.com/redis/go-redis/v9"
    "github.com/vvatanabe/dockertestx"
)

func TestRedis(t *testing.T) {
    // Start a Redis container with default options
    client, cleanup := dockertestx.NewRedis(t)
    defer cleanup()

    ctx := context.Background()

    // Prepare test data
    items := map[string]interface{}{
        "key1": "value1",
        "key2": "value2",
    }

    // Set up test data using PrepRedis
    if err := dockertestx.PrepRedis(t, client, items, time.Hour); err != nil {
        t.Fatalf("PrepRedis failed: %v", err)
    }

    // Verify the data
    for key, value := range items {
        got, err := client.Get(ctx, key).Result()
        if err != nil {
            t.Fatalf("failed to get item '%s': %v", key, err)
        }
        if got != value {
            t.Errorf("expected value '%v' for key '%s', but got '%v'",
                value, key, got)
        }
    }

    // Example with Redis List
    listKey := "mylist"
    listValues := []interface{}{"item1", "item2", "item3"}
    if err := dockertestx.PrepRedisList(t, client, listKey, listValues); err != nil {
        t.Fatalf("PrepRedisList failed: %v", err)
    }

    // Example with Redis Hash
    hashKey := "myhash"
    hashFields := map[string]interface{}{
        "field1": "value1",
        "field2": "value2",
    }
    if err := dockertestx.PrepRedisHash(t, client, hashKey, hashFields); err != nil {
        t.Fatalf("PrepRedisHash failed: %v", err)
    }
}
```

### MinIO (S3) Example
```go
package dockertestx_test

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vvatanabe/dockertestx"
)

func TestMinIO(t *testing.T) {
	// Start a MinIO container with default options
	client, cleanup := dockertestx.NewMinIO(t)
	defer cleanup()

	// Define a test bucket name
	bucketName := "test-bucket"
	ctx := context.Background()

	// Prepare test data
	testObjects := map[string][]byte{
		"test-file-1.txt": []byte("Hello, MinIO!"),
		"test-file-2.txt": []byte("This is a test file"),
		"dir/test-file-3.txt": []byte("Nested file test"),
	}

	// Set up test bucket and objects
	err := dockertestx.PrepS3Objects(t, client, bucketName, testObjects)
	if err != nil {
		t.Fatalf("PrepS3Objects failed: %v", err)
	}

	// Verify objects were uploaded correctly
	for key, expectedContent := range testObjects {
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			t.Fatalf("Failed to get object '%s': %v", key, err)
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read object '%s': %v", key, err)
		}

		if string(data) != string(expectedContent) {
			t.Errorf("Object '%s' content mismatch. Expected '%s', got '%s'", 
				key, string(expectedContent), string(data))
		}
	}
}
```

## Running Tests

Since **dockertestx** is intended for use in unit tests, you can run your tests as usual:

```bash
go test -v ./...
```

## Acknowledgments

- [dockertest](https://github.com/ory/dockertest) helps you boot up ephermal docker images for your Go tests with minimal work.
- [dynamotest](https://github.com/upsidr/dynamotest) is a package to help set up a DynamoDB Local Docker instance on your machine as a part of Go test code.
- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) provides APIs and utilities used for interfacing with AWS services, used here for S3-compatible storage.

## Authors

* **[vvatanabe](https://github.com/vvatanabe/)** - *Main contributor*
* Currently, there are no other contributors

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
