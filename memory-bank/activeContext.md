# Active Context

## Current Development Focus

The **dockertestx** library currently supports the following services, with each service implemented in a separate file:

1. **MySQL** (Implemented in `dockertestx.go`) - Basic SQL database functionality  
2. **PostgreSQL** (Implemented in `dockertestx.go`) - Basic SQL database functionality  
3. **Redis** (Implemented in `redis.go`) - Caching and in-memory data store  
4. **Memcached** (Implemented in `memcached.go`) - Simple caching service  
5. **MinIO** (Implemented in `minio.go`) - S3-compatible object storage  
6. **DynamoDB** (Implemented in `dynamodb.go`) - NoSQL database  

Each service follows a consistent implementation pattern:

1. A basic constructor function (e.g., `NewMySQL`)  
2. A constructor function with options (e.g., `NewMySQLWithOptions`)  
3. Helper functions for preparing test data (e.g., `PrepDatabase`, `PrepRedis`)  

## Recent Changes

Key recent updates include:

1. **Added DynamoDB support** - Introduced NoSQL database testing using the DynamoDB Local container.  
2. **Added MinIO support** - Enabled testing for S3-compatible object storage.  
3. **Enhanced Redis functionality** - Expanded support beyond basic key-value operations to include Lists, Hashes, Sets, and Sorted Sets.  
4. **Stabilized core functionality** - Improved error handling and recovery mechanisms.  

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
   // Example of a consistent API design
   func NewXxx(t testing.TB) (*xxx.Client, func())
   func NewXxxWithOptions(t testing.TB, runOpts []RunOption, hostOpts ...func(*docker.HostConfig)) (*xxx.Client, func())
   func PrepXxx(t testing.TB, client *xxx.Client, data ...) error
   ```
   
2. **Container resource management** - Evaluating options to configure resource limits (memory, CPU) for containers to optimize host machine usage.  
   ```go
   // Example: Applying resource limits via RunOption
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
   // Example: Loading data from external files
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