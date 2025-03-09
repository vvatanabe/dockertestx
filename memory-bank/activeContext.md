# Active Context

## Current Development Focus

The **dockertestx** library has been restructured with a modular package-based architecture. Each service is now implemented in its own dedicated package:

1. **MySQL & PostgreSQL** (Implemented in `sql` package) - Basic SQL database functionality  
2. **Redis** (Implemented in `redis` package) - Caching and in-memory data store  
3. **Memcached** (Implemented in `memcached` package) - Simple caching service  
4. **MinIO** (Implemented in `minio` package) - S3-compatible object storage  
5. **DynamoDB** (Implemented in `dynamodb` package) - NoSQL database  
6. **RabbitMQ** (Implemented in `rabbitmq` package) - Message broker
7. **Internal utilities** (Implemented in `internal` package) - Shared helper functions

Each service package follows a consistent implementation pattern:

1. A basic constructor function (e.g., `sql.RunMySQL`)  
2. A constructor function with options (e.g., `sql.RunMySQLWithOptions`)  
3. Helper functions for preparing test data (e.g., `sql.PrepDatabase`, `redis.PrepRedis`)  

Using this modular approach, consumers can import only the packages they need:

```go
// Import only the SQL package
import "github.com/vvatanabe/dockertestx/sql"

// Use directly from the package
db, cleanup := sql.RunMySQL(t)
defer cleanup()
```

## Recent Changes

Key recent updates include:

1. **Major package restructuring** - Transitioned from a single package with multiple files to a modular architecture with dedicated packages for each service.
2. **Added DynamoDB support** - Introduced NoSQL database testing using the DynamoDB Local container.  
3. **Added MinIO support** - Enabled testing for S3-compatible object storage.  
4. **Enhanced Redis functionality** - Expanded support beyond basic key-value operations to include Lists, Hashes, Sets, and Sorted Sets.  
5. **Added RabbitMQ support** - Implemented message broker testing with queue, exchange, binding operations and messaging capabilities.
6. **Stabilized core functionality** - Improved error handling and recovery mechanisms.  

## Next Steps

Upcoming development plans include:

1. **Adding MongoDB support** - Implementing support for a NoSQL document database.  
2. **Enhancing documentation** - Expanding examples and best practices.  
3. **Performance optimization** - Reducing container startup time and optimizing resource usage.  
4. **Increasing test coverage** - Adding more comprehensive test cases.  
5. **Exploring new APIs** - Enhancing functionality based on user feedback.  

## Ongoing Decisions & Considerations

Currently under discussion or implementation:

1. **API consistency** - Ensuring new services follow existing design patterns to minimize learning costs.  
   ```go
   // Example of a consistent API design across packages
   package xxx
   
   func RunXxx(t testing.TB) (*xxx.Client, func())
   func RunXxxWithOptions(t testing.TB, runOpts []RunOption, hostOpts ...func(*docker.HostConfig)) (*xxx.Client, func())
   func PrepXxx(t testing.TB, client *xxx.Client, data ...) error
   ```
   
2. **Container resource management** - Evaluating options to configure resource limits (memory, CPU) for containers to optimize host machine usage.  
   ```go
   // Example: Applying resource limits via RunOption in a specific package
   package redis
   
   func WithMemoryLimit(limit string) RunOption {
       return func(r *dockertest.RunOptions) {
           r.HostConfig = &docker.HostConfig{
               Resources: docker.Resources{
                   Memory: parse(limit),
               },
           }
       }
   }
   ```
   
3. **Improving error handling** - Providing more detailed error information and recovery strategies.  
   ```go
   // Example: Error with contextual information
   if err != nil {
       return fmt.Errorf("failed to connect to container %s: %w", containerName, err)
   }
   ```
   
4. **Parallel test execution support** - Enhancing stability and mitigating resource conflicts in concurrent runs.  

5. **Standardizing initial data population patterns** - Establishing efficient methods to load large datasets and complex structures for testing.  
   ```go
   // Example: Loading data from external files in the sql package
   package sql
   
   func PrepDatabaseFromFile(t testing.TB, db *sql.DB, schemaFile, dataFile string) error {
       // Read and execute SQL from files
   }
   ```

## Non-Functional Considerations

1. **Performance** - Optimizing container startup and data preparation to minimize test execution time.  
2. **Memory consumption** - Managing resources efficiently when running multiple containers in large test suites.  
3. **Test stability** - Reducing environmental discrepancies that cause inconsistent test results.  
4. **Usability** - Enhancing developer experience through intuitive APIs and clear error messages.  
5. **Ecosystem integration** - Ensuring seamless compatibility with other Go testing tools.  

## Challenges & Risks

1. **Docker dependency** - Tests require Docker, which may necessitate special configurations in CI/CD environments.  
2. **Network port conflicts** - Addressing potential port conflicts when running tests concurrently.  
3. **Image compatibility** - Maintaining compatibility across different Docker image versions.  

## Communication Plan

1. **Documentation** - Providing detailed usage guidelines and examples via README and GoDoc.  
2. **Release notes** - Clearly outlining changes and compatibility information in version updates.  
3. **Sample code** - Offering example use cases to illustrate common scenarios.
