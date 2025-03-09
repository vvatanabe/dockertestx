# dockertestx

**dockertestx** is a Go testing library that leverages [dockertest](https://github.com/ory/dockertest) to simplify integration testing with databases and other services using Docker containers. Originally forked from [sqltest](https://github.com/vvatanabe/sqltest), it extends the functionality to support not only SQL databases (MySQL, PostgreSQL) but also other data stores and services, making it easier to set up and manage Docker containers for your integration tests.

## Features

### Supported Services
- **SQL Databases**: MySQL 8.0 and PostgreSQL 13 container management
- **Cache Services**: Redis 7.2 and Memcached 1.6.18 support
- **Object Storage**: MinIO (S3-compatible) support
- **NoSQL Databases**: DynamoDB Local support
- **Message Brokers**: RabbitMQ support
- **Future Support**: MongoDB, Kafka, and other data stores
- **Extensibility**: Easy to add custom service containers

### Simple & Powerful
- **Easy to Use**: Clean and intuitive API
- **Automatic**: Container lifecycle management and health checks
- **Flexible**: Rich customization options and configurations
- **Reliable**: Built-in test helpers and utilities
- **Modular**: Import only what you need, reducing dependencies

## Installation

The library uses a modular package structure. You can install:

```bash
go get github.com/vvatanabe/dockertestx
```

Import only specific packages as needed：

```go
import "github.com/vvatanabe/dockertestx/sql"
import "github.com/vvatanabe/dockertestx/redis"
import "github.com/vvatanabe/dockertestx/memcached"
import "github.com/vvatanabe/dockertestx/minio"
import "github.com/vvatanabe/dockertestx/dynamodb"
import "github.com/vvatanabe/dockertestx/rabbitmq"
```

## Usage

For detailed usage examples, refer to the test files in each package:

- **SQL Package**: See [sql/sql_test.go](https://github.com/vvatanabe/sqltest/blob/main/sql/sql_test.go) for MySQL and PostgreSQL examples
- **Redis Package**: See [redis/redis_test.go](https://github.com/vvatanabe/sqltest/blob/main/redis/redis_test.go) for Redis examples
- **Memcached Package**: See [memcached/memcached_test.go](https://github.com/vvatanabe/sqltest/blob/main/memcached/memcached_test.go) for Memcached examples
- **MinIO Package**: See [minio/minio_test.go](https://github.com/vvatanabe/sqltest/blob/main/minio/minio_test.go) for S3-compatible storage examples
- **DynamoDB Package**: See [dynamodb/dynamodb_test.go](https://github.com/vvatanabe/sqltest/blob/main/dynamodb/dynamodb_test.go) for DynamoDB examples
- **RabbitMQ Package**: See [rabbitmq/rabbitmq_test.go](https://github.com/vvatanabe/sqltest/blob/main/rabbitmq/rabbitmq_test.go) for RabbitMQ examples

These test files demonstrate how to start containers, establish connections, and prepare test data for each supported service.

## Running Tests

Since **dockertestx** is intended for use in unit tests, you can run your tests as usual:

```bash
go test -v ./...
```

## Acknowledgments

- [dockertest](https://github.com/ory/dockertest) helps you boot up ephermal docker images for your Go tests with minimal work.
- [dynamotest](https://github.com/upsidr/dynamotest) is a package to help set up a DynamoDB Local Docker instance on your machine as a part of Go test code.
- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) provides APIs and utilities used for interfacing with AWS services, used here for S3-compatible storage and DynamoDB Local.
- [streadway/amqp](https://github.com/streadway/amqp) Go client for AMQP 0.9.1, used for RabbitMQ integration.

## **Authors**  

- **[vvatanabe](https://github.com/vvatanabe/)** - *Navigator* 🚀
- **Cline (a.k.a. Ryomen Sukuna 👹)** - *Driver* 🖋️

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
