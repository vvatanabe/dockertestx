package rabbitmq_test

import (
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/streadway/amqp"
	rabbitmqtest "github.com/vvatanabe/dockertestx/rabbitmq"
)

// TestDefaultRabbitMQ demonstrates using Run with default options.
func TestDefaultRabbitMQ(t *testing.T) {
	// Start a RabbitMQ container with default options.
	conn, cleanup := rabbitmqtest.Run(t)
	defer cleanup()

	// Test creating a channel
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	t.Log("Successfully connected to RabbitMQ and created a channel")
}

// TestRabbitMQWithCustomRunOptions demonstrates overriding default RunOptions.
func TestRabbitMQWithCustomRunOptions(t *testing.T) {
	// Custom RunOption to override the default tag
	customTag := func(opts *dockertest.RunOptions) {
		opts.Tag = "3-alpine" // Use a lighter version
	}

	// Start a RabbitMQ container with a custom tag
	conn, cleanup := rabbitmqtest.RunWithOptions(t, []func(*dockertest.RunOptions){customTag})
	defer cleanup()

	// Test basic functionality
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	t.Log("Successfully connected to RabbitMQ (alpine) and created a channel")
}

// TestRabbitMQWithCustomHostOptions demonstrates providing host configuration options.
func TestRabbitMQWithCustomHostOptions(t *testing.T) {
	// Host option to set AutoRemove to true
	autoRemove := func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	}

	// Start a RabbitMQ container with AutoRemove option
	conn, cleanup := rabbitmqtest.RunWithOptions(t, nil, autoRemove)
	defer cleanup()

	// Test basic functionality
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	t.Log("Successfully connected to RabbitMQ (with AutoRemove) and created a channel")
}

// TestRabbitMQQueueOperations tests queue creation and basic operations.
func TestRabbitMQQueueOperations(t *testing.T) {
	// Start a RabbitMQ container
	conn, cleanup := rabbitmqtest.Run(t)
	defer cleanup()

	// Create a queue
	queueName := "test-queue"
	options := amqp.Table{
		"durable": true,
	}

	queue, err := rabbitmqtest.PrepQueue(t, conn, queueName, options)
	if err != nil {
		t.Fatalf("failed to create queue: %v", err)
	}

	if queue.Name != queueName {
		t.Errorf("expected queue name '%s', got '%s'", queueName, queue.Name)
	}

	// Verify queue exists by declaring it again (should not error)
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclarePassive(
		queueName,
		true,  // durable (same as in the options above)
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		t.Fatalf("queue does not exist: %v", err)
	}

	t.Log("Successfully created and verified queue existence")
}

// TestRabbitMQExchangeOperations tests exchange creation.
func TestRabbitMQExchangeOperations(t *testing.T) {
	// Start a RabbitMQ container
	conn, cleanup := rabbitmqtest.Run(t)
	defer cleanup()

	// Create an exchange
	exchangeName := "test-exchange"
	exchangeType := "direct"
	options := amqp.Table{
		"durable": true,
	}

	err := rabbitmqtest.PrepExchange(t, conn, exchangeName, exchangeType, options)
	if err != nil {
		t.Fatalf("failed to create exchange: %v", err)
	}

	// Verify exchange exists by declaring it again (should not error)
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclarePassive(
		exchangeName,
		exchangeType,
		true,  // durable (same as in the options above)
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		t.Fatalf("exchange does not exist: %v", err)
	}

	t.Log("Successfully created and verified exchange existence")
}

// TestRabbitMQBindingOperations tests binding a queue to an exchange.
func TestRabbitMQBindingOperations(t *testing.T) {
	// Start a RabbitMQ container
	conn, cleanup := rabbitmqtest.Run(t)
	defer cleanup()

	// Create a queue
	queueName := "test-queue-binding"
	queueOptions := amqp.Table{
		"durable": true,
	}

	_, err := rabbitmqtest.PrepQueue(t, conn, queueName, queueOptions)
	if err != nil {
		t.Fatalf("failed to create queue: %v", err)
	}

	// Create an exchange
	exchangeName := "test-exchange-binding"
	exchangeType := "direct"
	exchangeOptions := amqp.Table{
		"durable": true,
	}

	err = rabbitmqtest.PrepExchange(t, conn, exchangeName, exchangeType, exchangeOptions)
	if err != nil {
		t.Fatalf("failed to create exchange: %v", err)
	}

	// Create a binding
	routingKey := "test-routing-key"
	bindingOptions := amqp.Table{}

	err = rabbitmqtest.PrepBinding(t, conn, queueName, exchangeName, routingKey, bindingOptions)
	if err != nil {
		t.Fatalf("failed to create binding: %v", err)
	}

	t.Log("Successfully created a binding between queue and exchange")
}

// TestRabbitMQPublishConsume tests publishing and consuming messages.
func TestRabbitMQPublishConsume(t *testing.T) {
	// Start a RabbitMQ container
	conn, cleanup := rabbitmqtest.Run(t)
	defer cleanup()

	// Create a queue
	queueName := "test-queue-messaging"
	_, err := rabbitmqtest.PrepQueue(t, conn, queueName, nil)
	if err != nil {
		t.Fatalf("failed to create queue: %v", err)
	}

	// Set up a consumer
	deliveries, consumerCleanup, err := rabbitmqtest.ConsumeMessages(t, conn, queueName)
	if err != nil {
		t.Fatalf("failed to set up consumer: %v", err)
	}
	defer consumerCleanup()

	// Publish a message
	message := []byte("Hello, RabbitMQ!")
	publishOptions := amqp.Publishing{
		ContentType: "text/plain",
	}

	err = rabbitmqtest.PublishMessage(t, conn, "", queueName, message, publishOptions)
	if err != nil {
		t.Fatalf("failed to publish message: %v", err)
	}

	// Wait for the message
	select {
	case delivery := <-deliveries:
		// Check message content
		if string(delivery.Body) != string(message) {
			t.Errorf("expected message '%s', got '%s'", string(message), string(delivery.Body))
		}
		// Acknowledge the message
		if err := delivery.Ack(false); err != nil {
			t.Errorf("failed to acknowledge message: %v", err)
		}
		t.Log("Successfully received and acknowledged message")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for message")
	}
}
