package dockertestx_test

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/vvatanabe/dockertestx"
)

func TestMinIO(t *testing.T) {
	// Start a MinIO container with default options
	client, cleanup := dockertestx.NewMinIO(t)
	defer cleanup()

	// Define a test bucket name
	bucketName := "test-bucket"
	ctx := context.Background()

	// Create the test bucket
	err := dockertestx.PrepBucket(t, client, bucketName)
	if err != nil {
		t.Fatalf("PrepBucket failed: %v", err)
	}

	// List buckets to verify our bucket was created
	buckets, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		t.Fatalf("Failed to list buckets: %v", err)
	}

	bucketFound := false
	for _, bucket := range buckets.Buckets {
		if *bucket.Name == bucketName {
			bucketFound = true
			break
		}
	}

	if !bucketFound {
		t.Fatalf("Created bucket '%s' not found", bucketName)
	}

	// Test uploading objects to the bucket
	testObjects := map[string][]byte{
		"test-file-1.txt":     []byte("Hello, MinIO!"),
		"test-file-2.txt":     []byte("This is a test file"),
		"dir/test-file-3.txt": []byte("Nested file test"),
	}

	err = dockertestx.PrepS3Objects(t, client, bucketName, testObjects)
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

func TestMinIOWithOptions(t *testing.T) {
	// Define custom options for the MinIO container
	runOpts := []dockertestx.RunOption{
		func(o *dockertest.RunOptions) {
			o.Tag = "RELEASE.2023-05-04T21-44-30Z"     // Specific MinIO version
			o.Env = append(o.Env, "MINIO_BROWSER=off") // Disable web UI
		},
	}

	// Apply custom host config options
	hostConfigOpts := []func(*docker.HostConfig){
		func(h *docker.HostConfig) {
			h.AutoRemove = true
		},
	}

	// Start MinIO with custom options
	client, cleanup := dockertestx.NewMinIOWithOptions(t, runOpts, hostConfigOpts...)
	defer cleanup()

	// Test that the client works by creating a bucket
	bucketName := "custom-options-test"
	err := dockertestx.PrepBucket(t, client, bucketName)
	if err != nil {
		t.Fatalf("PrepBucket failed with custom options: %v", err)
	}

	// Verify bucket was created
	ctx := context.Background()
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to find bucket '%s': %v", bucketName, err)
	}
}
