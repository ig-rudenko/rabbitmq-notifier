package producer

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"multiple-notifier/pkg/rabbitmq"
	"time"
)

type Config struct {
	ExchangeName string
	ExchangeType string
}

func NewProducer(config Config, rabbit *rabbitmq.Rabbit) *Producer {
	return &Producer{config, rabbit}
}

type Producer struct {
	config Config
	Rabbit *rabbitmq.Rabbit
}

func (p *Producer) Send(routingKey, data string) error {
	con, err := p.Rabbit.Connection()
	if err != nil {
		return err
	}

	chn, err := con.Channel()
	if err != nil {
		return err
	}

	if err := chn.ExchangeDeclare(
		p.config.ExchangeName,
		p.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         []byte(data),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = chn.PublishWithContext(ctx, p.config.ExchangeName, routingKey, false, false, msg)
	if err != nil {
		return err
	}

	return nil

}
