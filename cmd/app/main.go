package main

import (
	"log"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/pkg/rabbitmq"
	"time"
)

func main() {
	// RabbitMQ
	rc := rabbitmq.RabbitConfig{
		Schema:         "amqp",
		Username:       "rmuser",
		Password:       "rmpassword",
		Host:           "10.29.29.33",
		Port:           "5671",
		VHost:          "",
		ConnectionName: "my_app_name",
	}
	rbt := rabbitmq.NewRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("unable to connect to rabbit", err)
	}
	//

	// Consumer
	cc := consumer.ConsumerConfig{
		ExchangeName:  "user",
		ExchangeType:  "direct",
		RoutingKey:    "create",
		QueueName:     "test",
		ConsumerName:  "my_app_name",
		ConsumerCount: 3,
		PrefetchCount: 1,
	}
	cc.Reconnect.MaxAttempt = 60
	cc.Reconnect.Interval = 1 * time.Second
	csm := consumer.NewConsumer(cc, rbt)
	if err := csm.Start(); err != nil {
		log.Fatalln("unable to start consumer", err)
	}
	//

	select {}
}
