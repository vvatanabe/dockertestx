# dockertestx

**dockertestx** is a Go testing library that leverages [dockertest](https://github.com/ory/dockertest) to simplify integration testing with databases and other services using Docker containers. Originally forked from [sqltest](https://github.com/vvatanabe/sqltest), it extends the functionality to support not only SQL databases (MySQL, PostgreSQL) but also other data stores and services, making it easier to set up and manage Docker containers for your integration tests.

## Features

### Supported Services
- **SQL Databases**: MySQL and PostgreSQL container management
- **Cache Services**: Memcached support
- **Future Support**: MongoDB, Redis, and other data stores
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

### NewDockerDB
A helper function that starts a Docker container with the given run options, waits for the database to be ready, and returns a connected `*sql.DB` along with a cleanup function.

### PrepDatabase
Prepares the test database by executing provided schema (DDL) and initial data (DML) SQL statements. The initial data insertion is performed within a transaction to ensure consistency.

### InitialDBSetup
A helper struct used with `PrepDatabase` to specify the schema and initial data for setting up your test database.

## Examples

### MySQL Example
[Previous MySQL example code remains the same]

### PostgreSQL Example
[Previous PostgreSQL example code remains the same]

### Memcached Example

```go
package dockertestx_test

import (
    "testing"
    "github.com/bradfitz/gomemcache/memcache"
    "github.com/vvatanabe/dockertestx"
)

func TestMemcached(t *testing.T) {
    // Start a Memcached container with default options
    client, cleanup := dockertestx.NewMemcached(t)
    defer cleanup()

    // Prepare test data
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

    // Set up test data using PrepMemcached
    if err := dockertestx.PrepMemcached(t, client, items...); err != nil {
        t.Fatalf("PrepMemcached failed: %v", err)
    }

    // Verify the data
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
```

## Running Tests

Since **dockertestx** is intended for use in unit tests, you can run your tests as usual:

```bash
go test -v ./...
```

## Acknowledgments

- [dockertest](https://github.com/ory/dockertest) helps you boot up ephermal docker images for your Go tests with minimal work.
- [dynamotest](https://github.com/upsidr/dynamotest) is a package to help set up a DynamoDB Local Docker instance on your machine as a part of Go test code.

## Authors

* **[vvatanabe](https://github.com/vvatanabe/)** - *Main contributor*
* Currently, there are no other contributors

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
