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

- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)** – MySQL driver for establishing connections.  
  - Connects to MySQL containers.  
  - Executes SQL queries.  

- **[lib/pq](https://github.com/lib/pq)** – PostgreSQL driver for managing database connections.  
  - Connects to PostgreSQL containers.  
  - Executes SQL queries.  

- **[redis/go-redis](https://github.com/redis/go-redis) v9** – Redis client library.  
  - Establishes connections with Redis containers.  
  - Executes Redis commands.  

- **[bradfitz/gomemcache](https://github.com/bradfitz/gomemcache)** – Memcached client.  
  - Connects to Memcached containers.  
  - Executes Memcached commands.  

- **[AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)** – SDK for interacting with AWS-compatible services.  
  - Connects to S3 (MinIO).  
  - Connects to DynamoDB.  
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
┌─────────────────────┐
│    Go Application   │
└──────────┬──────────┘
           │
┌──────────v──────────┐
│     dockertestx     │
└──────────┬──────────┘
           │
┌──────────v──────────┐
│      dockertest     │
└──────────┬──────────┘
           │
┌──────────v──────────┐
│  Docker Engine API  │
└──────────┬──────────┘
           │
┌──────────v──────────┐┌─────────┐┌─────────┐┌─────────┐┌──────────┐
│        MySQL        ││Postgres ││  Redis  ││Memcached││ DynamoDB │
└─────────────────────┘└─────────┘└─────────┘└─────────┘└──────────┘
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
   go get github.com/vvatanabe/dockertestx
   ```

## Testing

To run the package tests, use the following commands:

```bash
# Run all tests
go test ./...

# Run specific tests (e.g., Redis-related tests only)
go test -run TestRedis

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
   db, cleanup := dockertestx.NewMySQL(t)
   defer cleanup() // Ensure resources are released
   ```

2. **Set Appropriate Timeouts**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

3. **Handle Errors Properly**
   ```go
   if err := dockertestx.PrepDatabase(t, db, setup); err != nil {
       t.Fatalf("failed to prepare database: %v", err)
   }
   ```

4. **Ensure Test Isolation**
   - Each test function should create its own container to avoid dependencies between tests.  

## Future Technical Roadmap

1. **Container Orchestration** – Enhancing multi-container interaction capabilities.  
2. **Support for New Data Stores** – Expanding support to MongoDB, Elasticsearch, etc.  
3. **Performance Optimizations** – Reducing container startup and teardown time.  
4. **Improved Parallel Test Execution** – Enhancing stability during concurrent test runs.  
5. **Cloud-Native Service Emulation** – Expanding support for AWS services and other cloud platforms.  
