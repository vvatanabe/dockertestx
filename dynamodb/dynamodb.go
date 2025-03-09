package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/vvatanabe/dockertestx"
	"testing"
	"time"
)

const (
	defaultDynamoDBImage = "amazon/dynamodb-local"
	defaultDynamoDBTag   = "latest"
	defaultRegion        = "us-east-1"
)

// NewDynamoDB starts a DynamoDB Local Docker container using the default settings and returns
// a connected *dynamodb.Client along with a cleanup function. It uses the DynamoDB Local image
// ("amazon/dynamodb-local") with tag "latest". For more customization, use NewDynamoDBWithOptions.
func NewDynamoDB(t testing.TB) (*dynamodb.Client, func()) {
	return NewDynamoDBWithOptions(t, nil)
}

// NewDynamoDBWithOptions starts a DynamoDB Local Docker container and returns a connected
// *dynamodb.Client along with a cleanup function. It applies the default settings:
//   - Repository: "amazon/dynamodb-local"
//   - Tag: "latest"
//
// Additional RunOption functions can be provided via the runOpts parameter to override these defaults,
// and optional host configuration functions can be provided via hostOpts.
func NewDynamoDBWithOptions(t testing.TB, runOpts []dockertestx.RunOption, hostOpts ...func(*docker.HostConfig)) (*dynamodb.Client, func()) {
	t.Helper()

	// Set default options for DynamoDB Local
	defaultRunOpts := &dockertest.RunOptions{
		Repository: defaultDynamoDBImage,
		Tag:        defaultDynamoDBTag,
	}

	// Apply any provided RunOption functions to override defaults
	for _, opt := range runOpts {
		opt(defaultRunOpts)
	}

	// Create a new Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to connect to docker: %s", err)
	}

	// Start the container with options
	resource, err := pool.RunWithOptions(defaultRunOpts, hostOpts...)
	if err != nil {
		t.Fatalf("failed to start dynamodb container: %s", err)
	}

	// Get the mapped port
	actualPort := resource.GetPort("8000/tcp")
	if actualPort == "" {
		_ = pool.Purge(resource)
		t.Fatalf("no host port was assigned for the dynamodb container")
	}
	t.Logf("DynamoDB container is running on host port '%s'", actualPort)

	// Configure AWS SDK v2
	endpoint := fmt.Sprintf("http://localhost:%s", actualPort)

	// Create a DynamoDB client with retry mechanism
	var client *dynamodb.Client
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = pool.Retry(func() error {
		// Configure AWS SDK credentials and endpoint
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: defaultRegion,
			}, nil
		})

		// Create AWS config
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(defaultRegion),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		)
		if err != nil {
			return fmt.Errorf("failed to configure AWS SDK: %w", err)
		}

		// Create DynamoDB client
		client = dynamodb.NewFromConfig(cfg)

		// Test connection with a simple ListTables call
		_, err = client.ListTables(ctx, &dynamodb.ListTablesInput{
			Limit: aws.Int32(1),
		})
		return err
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("failed to connect to dynamodb: %s", err)
	}

	// Create cleanup function
	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Logf("failed to remove dynamodb container: %s", err)
		}
	}

	return client, cleanup
}

// CreateDynamoDBTable creates a DynamoDB table with the given name, key schema, and attribute definitions.
// If the table already exists, it will not return an error.
func CreateDynamoDBTable(t testing.TB, client *dynamodb.Client, tableName string, keySchema []types.KeySchemaElement, attributeDefs []types.AttributeDefinition) error {
	t.Helper()

	ctx := context.Background()

	// Check if table already exists
	tables, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	for _, existingTable := range tables.TableNames {
		if existingTable == tableName {
			t.Logf("Table %s already exists", tableName)
			return nil
		}
	}

	// Create table
	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		KeySchema:            keySchema,
		AttributeDefinitions: attributeDefs,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	t.Logf("Created table %s", tableName)
	return nil
}

// PrepDynamoDBItems inserts the provided items into the specified DynamoDB table.
// It accepts tableName and a list of items as map[string]types.AttributeValue.
func PrepDynamoDBItems(t testing.TB, client *dynamodb.Client, tableName string, items []map[string]types.AttributeValue) error {
	t.Helper()

	ctx := context.Background()

	for i, item := range items {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})

		if err != nil {
			return fmt.Errorf("failed to insert item %d into table %s: %w", i, tableName, err)
		}
	}

	t.Logf("Inserted %d items into table %s", len(items), tableName)
	return nil
}

// DeleteDynamoDBTable deletes the specified DynamoDB table.
// It's useful for cleanup after tests.
func DeleteDynamoDBTable(t testing.TB, client *dynamodb.Client, tableName string) error {
	t.Helper()

	ctx := context.Background()

	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return fmt.Errorf("failed to delete table %s: %w", tableName, err)
	}

	t.Logf("Deleted table %s", tableName)
	return nil
}
