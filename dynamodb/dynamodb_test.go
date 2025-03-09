package dynamodb_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ory/dockertest/v3"
	"github.com/vvatanabe/dockertestx"
	dynamodbtest "github.com/vvatanabe/dockertestx/dynamodb"
	"testing"
)

func TestDynamoDB(t *testing.T) {
	// Start a DynamoDB container with default options
	client, cleanup := dynamodbtest.NewDynamoDB(t)
	defer cleanup()

	ctx := context.Background()

	// Define table structure
	tableName := "TestUsers"
	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("ID"),
			KeyType:       types.KeyTypeHash,
		},
	}
	attrDefs := []types.AttributeDefinition{
		{
			AttributeName: aws.String("ID"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	// Create table
	err := dynamodbtest.CreateDynamoDBTable(t, client, tableName, keySchema, attrDefs)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Define test items
	type User struct {
		ID   string `dynamodbav:"ID"`
		Name string `dynamodbav:"Name"`
		Age  int    `dynamodbav:"Age"`
	}

	users := []User{
		{ID: "1", Name: "Alice", Age: 30},
		{ID: "2", Name: "Bob", Age: 25},
		{ID: "3", Name: "Charlie", Age: 35},
	}

	// Marshal and insert items
	var items []map[string]types.AttributeValue
	for _, user := range users {
		item, err := attributevalue.MarshalMap(user)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}
		items = append(items, item)
	}

	err = dynamodbtest.PrepDynamoDBItems(t, client, tableName, items)
	if err != nil {
		t.Fatalf("Failed to insert items: %v", err)
	}

	// Query items to verify
	for _, user := range users {
		key, err := attributevalue.MarshalMap(map[string]string{
			"ID": user.ID,
		})
		if err != nil {
			t.Fatalf("Failed to marshal key: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       key,
		})
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}

		if resp.Item == nil {
			t.Fatalf("Item not found: %s", user.ID)
		}

		var fetchedUser User
		err = attributevalue.UnmarshalMap(resp.Item, &fetchedUser)
		if err != nil {
			t.Fatalf("Failed to unmarshal item: %v", err)
		}

		if fetchedUser.ID != user.ID || fetchedUser.Name != user.Name || fetchedUser.Age != user.Age {
			t.Errorf("Item mismatch for ID %s. Expected %+v, got %+v", user.ID, user, fetchedUser)
		}
	}

	// Test table deletion
	err = dynamodbtest.DeleteDynamoDBTable(t, client, tableName)
	if err != nil {
		t.Fatalf("Failed to delete table: %v", err)
	}
}

// TestDynamoDBWithOptions demonstrates how to customize the DynamoDB container
func TestDynamoDBWithOptions(t *testing.T) {
	// Use custom options
	runOpts := []dockertestx.RunOption{
		func(opts *dockertest.RunOptions) {
			opts.Tag = "1.21.0" // Specific version
			opts.Env = append(opts.Env, "AWS_ACCESS_KEY_ID=customkey")
			opts.Env = append(opts.Env, "AWS_SECRET_ACCESS_KEY=customsecret")
			// Add additional customization as needed
		},
	}

	// Start container with custom options
	client, cleanup := dynamodbtest.NewDynamoDBWithOptions(t, runOpts)
	defer cleanup()

	// Verify container works
	ctx := context.Background()
	_, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}
}
