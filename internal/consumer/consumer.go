package consumer

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"multiple-notifier/pkg/rabbitmq"
	"os"
	"time"
)

type Config struct {
	ExchangeName  string
	ExchangeType  string
	RoutingKey    string
	QueueName     string
	ConsumerName  string
	ConsumerCount int
	PrefetchCount int
	Reconnect     struct {
		MaxAttempt int
		Interval   time.Duration
	}
}

type Consumer struct {
	config Config
	Rabbit *rabbitmq.Rabbit
}

// NewConsumer returns a consumer instance.
func NewConsumer(config Config, rabbit *rabbitmq.Rabbit) *Consumer {
	return &Consumer{
		config: config,
		Rabbit: rabbit,
	}
}

// Start declares all the necessary components of the consumer and
// runs the consumers. This is called one at the application start up
// or when consumer needs to reconnects to the server.
func (c *Consumer) Start() error {
	con, err := c.Rabbit.Connection()
	if err != nil {
		return err
	}
	go c.closedConnectionListener(con.NotifyClose(make(chan *amqp.Error)))

	chn, err := con.Channel()
	if err != nil {
		return err
	}

	if err := chn.ExchangeDeclare(
		c.config.ExchangeName,
		c.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := chn.QueueDeclare(
		c.config.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-queue-mode": "lazy"},
	); err != nil {
		return err
	}

	if err := chn.QueueBind(
		c.config.QueueName,
		c.config.RoutingKey,
		c.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := chn.Qos(c.config.PrefetchCount, 0, false); err != nil {
		return err
	}

	for i := 1; i <= c.config.ConsumerCount; i++ {
		id := i
		go c.consume(chn, id)
	}

	// Simulate manual connection close
	//_ = con.Close()

	return nil
}

// closedConnectionListener attempts to reconnect to the server and
// reopens the channel for set amount of time if the connection is
// closed unexpectedly. The attempts are spaced at equal intervals.
func (c *Consumer) closedConnectionListener(closed <-chan *amqp.Error) {
	log.Println("INFO: Watching closed connection")

	// If you do not want to reconnect in the case of manual disconnection
	// via RabbitMQ UI or Server restart, handle `amqp.ConnectionForced`
	// error code.
	err := <-closed
	if err != nil {
		log.Println("INFO: Closed connection:", err.Error())

		var i int

		for i = 0; i < c.config.Reconnect.MaxAttempt; i++ {
			log.Println("INFO: Attempting to reconnect")

			if err := c.Rabbit.Connect(); err == nil {
				log.Println("INFO: Reconnected")

				if err := c.Start(); err == nil {
					break
				}
			}

			time.Sleep(c.config.Reconnect.Interval)
		}

		if i == c.config.Reconnect.MaxAttempt {
			log.Println("CRITICAL: Giving up reconnecting")

			return
		}
	} else {
		log.Println("INFO: Connection closed normally, will not reconnect")
		os.Exit(0)
	}
}

// consume creates a new consumer and starts consuming the messages.
// If this is called more than once, there will be multiple consumers
// running. All consumers operate in a round robin fashion to distribute
// message load.
func (c *Consumer) consume(channel *amqp.Channel, id int) {
	consumerName := fmt.Sprintf("%s (%d/%d)", c.config.ConsumerName, id, c.config.ConsumerCount)
	msgs, err := channel.Consume(
		c.config.QueueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(fmt.Sprintf("CRITICAL: Unable to start consumer (%d/%d)", id, c.config.ConsumerCount))

		return
	}

	log.Println("[", consumerName, "] Running ...")
	log.Println("[", consumerName, "] Press CTRL+C to exit ...")

	for msg := range msgs {
		log.Println("[", consumerName, "] Consumed:", string(msg.Body))

		if err := msg.Ack(false); err != nil {

			log.Println("Unable to acknowledge the message, dropped", err)
		}
	}

	log.Println("[", consumerName, "] Exiting ...")
}
