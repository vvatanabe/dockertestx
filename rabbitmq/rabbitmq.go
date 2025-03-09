package rabbitmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/streadway/amqp"
)

const (
	defaultRabbitMQImage = "rabbitmq"
	defaultRabbitMQTag   = "3-management"
	defaultAMQPPort      = "5672/tcp"
)

// Run starts a RabbitMQ Docker container using the default settings and returns a connected
// *amqp.Connection along with a cleanup function. It uses the default RabbitMQ image ("rabbitmq")
// with tag "3-management". For more customization, use RunWithOptions.
func Run(t testing.TB) (*amqp.Connection, func()) {
	return RunWithOptions(t, nil)
}

// RunWithOptions starts a RabbitMQ Docker container using Docker and returns a connected
// *amqp.Connection along with a cleanup function. It applies the default settings:
//   - Repository: "rabbitmq"
//   - Tag: "3-management"
//
// Additional RunOption functions can be provided via the runOpts parameter to override these defaults,
// and optional host configuration functions can be provided via hostOpts.
func RunWithOptions(t testing.TB, runOpts []func(*dockertest.RunOptions), hostOpts ...func(*docker.HostConfig)) (*amqp.Connection, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to connect to docker: %s", err)
	}

	// Set default run options for RabbitMQ
	defaultRunOpts := &dockertest.RunOptions{
		Repository: defaultRabbitMQImage,
		Tag:        defaultRabbitMQTag,
		Env: []string{
			"RABBITMQ_DEFAULT_USER=guest",
			"RABBITMQ_DEFAULT_PASS=guest",
		},
	}

	// Apply any provided RunOption functions to override defaults
	for _, opt := range runOpts {
		opt(defaultRunOpts)
	}

	// Pass optional host configuration options
	resource, err := pool.RunWithOptions(defaultRunOpts, hostOpts...)
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %s", err)
	}

	actualPort := resource.GetHostPort(defaultAMQPPort)
	if actualPort == "" {
		_ = pool.Purge(resource)
		t.Fatal("no host port was assigned for the rabbitmq container")
	}
	t.Logf("rabbitmq container is running on host port '%s'", actualPort)

	// Create RabbitMQ connection
	var conn *amqp.Connection

	// Try to connect to RabbitMQ with retries
	if err = pool.Retry(func() error {
		var err error
		conn, err = amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s/", actualPort))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("could not connect to rabbitmq: %s", err)
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close RabbitMQ connection: %s", err)
		}
		if err := pool.Purge(resource); err != nil {
			t.Logf("failed to remove rabbitmq container: %s", err)
		}
	}

	return conn, cleanup
}

// PrepQueue creates a queue in RabbitMQ with the specified name and options.
// It returns the created queue and an error if the operation fails.
func PrepQueue(t testing.TB, conn *amqp.Connection, name string, options amqp.Table) (*amqp.Queue, error) {
	t.Helper()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	if options != nil {
		if val, ok := options["durable"]; ok {
			if durableBool, ok := val.(bool); ok {
				durable = durableBool
			}
		}
		if val, ok := options["autoDelete"]; ok {
			if autoDeleteBool, ok := val.(bool); ok {
				autoDelete = autoDeleteBool
			}
		}
		if val, ok := options["exclusive"]; ok {
			if exclusiveBool, ok := val.(bool); ok {
				exclusive = exclusiveBool
			}
		}
		if val, ok := options["noWait"]; ok {
			if noWaitBool, ok := val.(bool); ok {
				noWait = noWaitBool
			}
		}
	}

	q, err := ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		options,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue '%s': %w", name, err)
	}

	return &q, nil
}

// PrepExchange creates an exchange in RabbitMQ with the specified name, kind, and options.
// It returns an error if the operation fails.
func PrepExchange(t testing.TB, conn *amqp.Connection, name string, kind string, options amqp.Table) error {
	t.Helper()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	durable := false
	autoDelete := false
	internal := false
	noWait := false

	if options != nil {
		if val, ok := options["durable"]; ok {
			if durableBool, ok := val.(bool); ok {
				durable = durableBool
			}
		}
		if val, ok := options["autoDelete"]; ok {
			if autoDeleteBool, ok := val.(bool); ok {
				autoDelete = autoDeleteBool
			}
		}
		if val, ok := options["internal"]; ok {
			if internalBool, ok := val.(bool); ok {
				internal = internalBool
			}
		}
		if val, ok := options["noWait"]; ok {
			if noWaitBool, ok := val.(bool); ok {
				noWait = noWaitBool
			}
		}
	}

	err = ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange '%s': %w", name, err)
	}

	return nil
}

// PrepBinding creates a binding between a queue and an exchange with the specified routing key and options.
// It returns an error if the operation fails.
func PrepBinding(t testing.TB, conn *amqp.Connection, queueName string, exchangeName string, routingKey string, options amqp.Table) error {
	t.Helper()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	err = ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // noWait
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue '%s' to exchange '%s': %w", queueName, exchangeName, err)
	}

	return nil
}

// PublishMessage publishes a message to the specified exchange with the routing key and options.
// It returns an error if the operation fails.
func PublishMessage(t testing.TB, conn *amqp.Connection, exchange string, routingKey string, message []byte, options amqp.Publishing) error {
	t.Helper()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Set default content type if not provided
	if options.ContentType == "" {
		options.ContentType = "text/plain"
	}

	// Set the message body
	options.Body = message

	err = ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to exchange '%s': %w", exchange, err)
	}

	return nil
}

// ConsumeMessages sets up a consumer for a queue and returns a channel for receiving messages.
// It also returns a function to cancel the consumer.
func ConsumeMessages(t testing.TB, conn *amqp.Connection, queueName string) (<-chan amqp.Delivery, func(), error) {
	t.Helper()

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	consumerName := fmt.Sprintf("consumer-%d", time.Now().UnixNano())
	deliveries, err := ch.Consume(
		queueName,
		consumerName,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("failed to register consumer for queue '%s': %w", queueName, err)
	}

	cleanup := func() {
		if err := ch.Cancel(consumerName, false); err != nil {
			t.Logf("failed to cancel consumer: %s", err)
		}
		if err := ch.Close(); err != nil {
			t.Logf("failed to close channel: %s", err)
		}
	}

	return deliveries, cleanup, nil
}
