package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/vvatanabe/dockertestx/internal"
	"strings"
	"testing"
	"time"
)

const (
	defaultMinIOImage = "minio/minio"
	defaultMinIOTag   = "latest"
	defaultAccessKey  = "minioadmin"
	defaultSecretKey  = "minioadmin"
)

// Run starts a MinIO Docker container using the default settings and returns a configured S3 client
// along with a cleanup function. It uses the default MinIO image ("minio/minio") with tag "latest".
// For more customization, use RunWithOptions.
func Run(t testing.TB) (*s3.Client, func()) {
	return RunWithOptions(t, nil)
}

// RunWithOptions starts a MinIO Docker container using Docker and returns a configured S3 client
// along with a cleanup function. It applies the default settings:
//   - Repository: "minio/minio"
//   - Tag: "latest"
//   - Environment: MINIO_ROOT_USER=minioadmin, MINIO_ROOT_PASSWORD=minioadmin
//   - Command: ["server", "/data"]
//
// Additional RunOption functions can be provided via the runOpts parameter to override these defaults,
// and optional host configuration functions can be provided via hostOpts.
func RunWithOptions(t testing.TB, runOpts []func(*dockertest.RunOptions), hostOpts ...func(*docker.HostConfig)) (*s3.Client, func()) {
	t.Helper()

	// Set default run options for MinIO
	defaultRunOpts := &dockertest.RunOptions{
		Repository: defaultMinIOImage,
		Tag:        defaultMinIOTag,
		Env: []string{
			"MINIO_ROOT_USER=" + defaultAccessKey,
			"MINIO_ROOT_PASSWORD=" + defaultSecretKey,
		},
		Cmd: []string{"server", "/data"},
	}

	// Apply any provided RunOption functions to override defaults
	for _, opt := range runOpts {
		opt(defaultRunOpts)
	}

	// Start the container
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to connect to docker: %s", err)
	}

	// Pass optional host configuration options
	resource, err := pool.RunWithOptions(defaultRunOpts, hostOpts...)
	if err != nil {
		t.Fatalf("failed to start MinIO container: %s", err)
	}

	// Get the port that MinIO is running on
	actualPort := resource.GetHostPort("9000/tcp")
	if actualPort == "" {
		_ = pool.Purge(resource)
		t.Fatalf("no host port was assigned for the MinIO container")
	}

	t.Logf("MinIO container is running on host port '%s'", actualPort)
	// GetHostPort may return a format like "localhost:55250",
	// so remove the "localhost:" prefix if present
	actualPort = strings.TrimPrefix(actualPort, "localhost:")

	// Get access and secret keys from environment variables
	accessKey := internal.GetEnvValue(defaultRunOpts.Env, "MINIO_ROOT_USER")
	secretKey := internal.GetEnvValue(defaultRunOpts.Env, "MINIO_ROOT_PASSWORD")

	// Wait for MinIO to be ready
	var s3Client *s3.Client
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set timeout to 3 seconds
	defer cancel()

	// Build the endpoint URL with the correct localhost:port format
	endpoint := fmt.Sprintf("http://localhost:%s", actualPort)
	t.Logf("Connecting to MinIO endpoint: %s with credentials %s:%s", endpoint, accessKey, secretKey)

	if err = pool.Retry(func() error {
		// Load AWS SDK Go v2 configuration
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		// Create S3 client with direct option settings
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // MinIO requires path-style addressing
			o.BaseEndpoint = aws.String(endpoint)
		})

		// Test if the S3 API is responding
		_, err = s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			t.Logf("Failed to list buckets, retrying: %v", err)
			return err
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("failed to connect to MinIO: %s", err)
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Logf("failed to remove MinIO container: %s", err)
		}
	}

	return s3Client, cleanup
}

// PrepBucket creates a bucket if it doesn't exist
func PrepBucket(t testing.TB, client *s3.Client, bucketName string) error {
	t.Helper()
	ctx := context.Background()

	// Check if bucket exists
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		// Create bucket if it doesn't exist
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	return nil
}

// UploadObject uploads an object to a bucket
func UploadObject(t testing.TB, client *s3.Client, bucketName, key string, body []byte) error {
	t.Helper()
	ctx := context.Background()

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object %s to bucket %s: %w", key, bucketName, err)
	}

	return nil
}

// PrepS3Objects prepares a bucket with the given objects
func PrepS3Objects(t testing.TB, client *s3.Client, bucketName string, objects map[string][]byte) error {
	t.Helper()

	// Create bucket if it doesn't exist
	if err := PrepBucket(t, client, bucketName); err != nil {
		return err
	}

	// Upload objects
	for key, data := range objects {
		if err := UploadObject(t, client, bucketName, key, data); err != nil {
			return err
		}
	}

	return nil
}
