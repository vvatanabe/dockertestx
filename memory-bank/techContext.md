# Tech Context

## Technologies Used

### Programming Language
- **Go** – The entire package is implemented in Go, adhering to standard Go patterns and conventions.  
- **Go Modules** – Dependency management is handled using Go Modules.  

### Key Dependencies

- **[dockertest](https://github.com/ory/dockertest) v3** – Core library for Docker container management.  
  - Manages container lifecycle (start, stop, remove).  
  - Implements health checks and retry logic.  
  - Abstracts low-level Docker API calls.  
  - Used by all service packages.

- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)** – MySQL driver for establishing connections.  
  - Connects to MySQL containers.  
  - Executes SQL queries.  
  - Used in the `sql` package.

- **[lib/pq](https://github.com/lib/pq)** – PostgreSQL driver for managing database connections.  
  - Connects to PostgreSQL containers.  
  - Executes SQL queries.  
  - Used in the `sql` package.

- **[redis/go-redis](https://github.com/redis/go-redis) v9** – Redis client library.  
  - Establishes connections with Redis containers.  
  - Executes Redis commands.  
  - Used in the `redis` package.

- **[bradfitz/gomemcache](https://github.com/bradfitz/gomemcache)** – Memcached client.  
  - Connects to Memcached containers.  
  - Executes Memcached commands.  
  - Used in the `memcached` package.

- **[AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)** – SDK for interacting with AWS-compatible services.  
  - Connects to S3 (MinIO) - Used in the `minio` package.
  - Connects to DynamoDB - Used in the `dynamodb` package.
  - Manages AWS credentials.  

### Supported Service Technologies

#### Databases
- **MySQL 8.0** – Default version used.  
- **PostgreSQL 13** – Default version used.  

#### Caching Systems
- **Redis 7.2** – Default version used.  
- **Memcached 1.6.18** – Default version used.  

#### Cloud Service Compatibility
- **MinIO** – S3-compatible object storage.  
- **DynamoDB Local** – Local version of AWS DynamoDB.  

### Technology Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                        Go Application                           │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌──────────────────┬──────────────┼──────────────┬────────────────┐
│      sql         │     redis    │   memcached  │     minio      │
│   Package        │    Package   │    Package   │    Package     │
└──────────┬───────┘──────┬───────┘──────┬───────┘──────┬─────────┘
           │              │              │              │
           │         ┌────┴───────────┐  │         ┌────┴─────────┐
           │         │    dynamodb    │  │         │   internal   │
           │         │    Package     │  │         │   Package    │
           │         └────────────────┘  │         └──────────────┘
           │              │              │              │
           └──────────────┴──────────────┴──────────────┘
                                  │
┌─────────────────────────────────┼───────────────────────────────┐
│                            dockertest                           │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────┼───────────────────────────────┐
│                      Docker Engine API                          │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌──────────────┬──────────────────┼──────────────┬────────────────┐
│    MySQL     │    PostgreSQL    │    Redis     │   Memcached    │
└──────────────┘──────────────────┘──────────────┘────────────────┘
           ┌────────────────────────┐       ┌────────────────────┐
           │       DynamoDB         │       │       MinIO        │
           └────────────────────────┘       └────────────────────┘
```

## Development Environment

### Requirements
- **Go 1.18+** – To leverage modern Go features.  
- **Docker** – Required for running containers locally.  
- **Docker Compose** (optional) – Useful for multi-container testing.  

### Setup Instructions

1. **Install Go**
   ```bash
   # macOS
   brew install go
   
   # Linux
   # Use the package manager specific to your distribution
   ```

2. **Install Docker**
   ```bash
   # macOS
   brew install --cask docker
   
   # Linux
   # Follow installation instructions for your distribution
   ```

3. **Start the Docker Daemon**
   ```bash
   # macOS
   open -a Docker
   
   # Linux
   sudo systemctl start docker
   ```

4. **Install dockertestx**
   ```bash
   # Install the entire library
   go get github.com/vvatanabe/dockertestx/...
   
   # Or install only specific packages as needed
   go get github.com/vvatanabe/dockertestx/sql
   go get github.com/vvatanabe/dockertestx/redis
   ```

## Testing

To run the package tests, use the following commands:

```bash
# Run all tests
go test ./...

# Run specific tests (e.g., Redis-related tests only)
go test ./redis -run TestRedis

# Run tests with verbose output
go test -v ./...

# Run tests with race condition detection
go test -race ./...
```

## Technical Constraints

1. **Docker Dependency** – Cannot run in environments without Docker installed.  
2. **Network Ports** – Uses temporary ports during test execution, so conflicts may occur.  
3. **Resource Usage** – Running multiple containers simultaneously can impact memory and CPU resources.  
4. **Cleanup Issues** – If a test is forcibly terminated, orphaned containers may remain.  

## Best Practices

1. **Ensure Proper Resource Cleanup**
   ```go
   import "github.com/vvatanabe/dockertestx/sql"
   
   // Using the sql package
   db, cleanup := sql.RunMySQL(t)
   defer cleanup() // Ensure resources are released
   ```

2. **Set Appropriate Timeouts**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

3. **Handle Errors Properly**
   ```go
   import "github.com/vvatanabe/dockertestx/sql"
   
   // Using the sql package's PrepDatabase function
   if err := sql.PrepDatabase(t, db, setup); err != nil {
       t.Fatalf("failed to prepare database: %v", err)
   }
   ```

4. **Import Only What You Need**
   ```go
   // If you only need Redis functionality
   import "github.com/vvatanabe/dockertestx/redis"
   
   // Use redis package directly
   client, cleanup := redis.RunRedis(t)
   defer cleanup()
   ```

4. **Ensure Test Isolation**
   - Each test function should create its own container to avoid dependencies between tests.  

## Future Technical Roadmap

1. **Container Orchestration** – Enhancing multi-container interaction capabilities.  
2. **Support for New Data Stores** – Expanding support to MongoDB, Elasticsearch, etc.  
3. **Performance Optimizations** – Reducing container startup and teardown time.  
4. **Improved Parallel Test Execution** – Enhancing stability during concurrent test runs.  
5. **Cloud-Native Service Emulation** – Expanding support for AWS services and other cloud platforms.
