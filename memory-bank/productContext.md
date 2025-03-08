# Product Context

## Problems to Solve

In integration testing, verifying interactions with real databases and services often requires complex environment setup. This leads to several challenges:

1. **Environment Dependency** – Different configurations are required for each developer’s local setup and CI/CD pipelines.  
2. **Setup Complexity** – Building a test environment involving multiple services is time-consuming and labor-intensive.  
3. **Test Consistency** – Differences in environments can cause inconsistent test results.  
4. **Difficult Cleanup** – Restoring the environment to its original state after tests is not always straightforward.  
5. **Lack of a Unified Approach for Different Services** – Each service requires different testing methodologies.  

## How dockertestx Solves These Problems

dockertestx addresses these issues through the following approaches:

1. **Leveraging Docker Containers** – Provides a consistent, environment-independent test environment.  
2. **Simple API** – Enables easy container management without requiring deep Docker knowledge.  
3. **Automated Lifecycle Management** – Handles container startup, connection establishment, and shutdown automatically.  
4. **Consistent Patterns** – Allows writing tests for different services (MySQL, PostgreSQL, Redis, DynamoDB, etc.) using a unified approach.  
5. **Helper Functions** – Simplifies test data preparation and reduces boilerplate code.  

## User Experience Goals

dockertestx aims to provide developers with the following experience:

1. **Simplicity** – Set up a test environment with just a few lines of code.  
2. **Consistency** – Write tests for various services using a uniform pattern.  
3. **Self-Containment** – Create, test, and clean up the entire environment solely within the test code.  
4. **Productivity** – Reduce the time spent on test environment setup and focus on writing actual tests.  
5. **Robustness** – Ensure a reliable and reproducible testing environment.  

## Example Use Cases

### Basic Usage

```go
func TestMyService(t *testing.T) {
    // Start a MySQL container and establish a connection
    db, cleanup := dockertestx.NewMySQL(t)
    defer cleanup() // Automatically stop and remove the container after the test

    // Prepare schema and initial data for testing
    err := dockertestx.PrepDatabase(t, db, dockertestx.InitialDBSetup{
        SchemaSQL: "CREATE TABLE users (id INT, name VARCHAR(255))",
        InitialData: []string{
            "INSERT INTO users VALUES (1, 'Alice')",
            "INSERT INTO users VALUES (2, 'Bob')",
        },
    })

    // Execute the code under test
    svc := NewService(db)
    result := svc.GetUsers()

    // Assertions
    assert.Equal(t, 2, len(result))
}
```

### Complex Test Cases

Testing interactions between multiple services is also straightforward:

```go
func TestComplexService(t *testing.T) {
    // Start a MySQL container
    db, dbCleanup := dockertestx.NewMySQL(t)
    defer dbCleanup()

    // Start a Redis container
    redis, redisCleanup := dockertestx.NewRedis(t)
    defer redisCleanup()

    // Start a MinIO (S3-compatible) container
    s3Client, s3Cleanup := dockertestx.NewMinIO(t)
    defer s3Cleanup()

    // Prepare initial data for each service
    // ...

    // Initialize the service under test with multiple backends
    svc := NewComplexService(db, redis, s3Client)

    // Execute the test
    // ...
}
```

## Future Vision

Beyond the currently supported services, dockertestx aims to expand compatibility with more data stores and services. Continuous improvements in features and user experience will be prioritized, with the ultimate goal of making dockertestx the standard choice for integration testing in the Go ecosystem.
