package main

import (
	"log"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/internal/misc"
	"multiple-notifier/internal/notifier/telegram"
	"multiple-notifier/pkg/rabbitmq"
	"os"
	"time"
)

func main() {
	// RabbitMQ
	rc := rabbitmq.Config{
		Schema:         "amqps",
		Username:       misc.GetEnvOrPanic("RABBITMQ_USER"),
		Password:       misc.GetEnvOrPanic("RABBITMQ_PASS"),
		Host:           misc.GetEnvOrPanic("RABBITMQ_HOST"),
		Port:           misc.GetEnv("RABBITMQ_PORT", "5671"),
		VHost:          misc.GetEnv("RABBITMQ_VHOST", ""),
		ConnectionName: misc.GetEnvOrPanic("RABBITMQ_CONNECTION_NAME"),
	}
	rbt := rabbitmq.NewRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("unable to connect to rabbit", err)
	}

	// Consumer
	cc := consumer.Config{
		ExchangeName:  misc.GetEnvOrPanic("RABBITMQ_EXCHANGE_NAME"),
		ExchangeType:  misc.GetEnv("RABBITMQ_EXCHANGE_TYPE", "direct"),
		RoutingKey:    misc.GetEnvOrPanic("RABBITMQ_ROUTING_KEY"),
		QueueName:     misc.GetEnvOrPanic("RABBITMQ_QUEUE"),
		ConsumerName:  misc.GetEnvOrPanic("RABBITMQ_CONNECTION_NAME"),
		ConsumerCount: misc.GetIntEnv("CONSUMER_COUNT", 5),
		PrefetchCount: misc.GetIntEnv("PREFETCH_COUNT", 5),
	}
	cc.Reconnect.MaxAttempt = 60
	cc.Reconnect.Interval = 1 * time.Second

	cc.Notifier = getNotifier()

	csm := consumer.NewConsumer(cc, rbt)
	if err := csm.Start(); err != nil {
		log.Fatalln("unable to start consumer", err)
	}
	//

	select {}
}

func getNotifier() consumer.Notifier {
	if os.Args[1] == "telegram" {
		return telegram.NewNotifier(misc.GetEnvOrPanic("TELEGRAM_TOKEN"), misc.GetEnvOrPanic("TELEGRAM_CHAT_ID"))
	}
	return nil
}
