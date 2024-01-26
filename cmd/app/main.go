package main

import (
	"fmt"
	"log"
	"multiple-notifier/internal/config"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/internal/mode"
	"multiple-notifier/internal/producer"
	"multiple-notifier/pkg/rabbitmq"
	"os"
	"time"
)

func main() {
	app := mode.NewApp()
	app.ParseArgs()

	mainConfig := config.NewConfig()

	// RabbitMQ
	rc := rabbitmq.Config{
		Schema:     "amqps",
		Username:   mainConfig.Rabbitmq.User,
		Password:   mainConfig.Rabbitmq.Password,
		Host:       mainConfig.Rabbitmq.Host,
		Port:       mainConfig.Rabbitmq.Port,
		VHost:      mainConfig.Rabbitmq.Vhost,
		CaCertFile: mainConfig.Rabbitmq.CaCert,
		CertFile:   mainConfig.Rabbitmq.CertFile,
		KeyFile:    mainConfig.Rabbitmq.KeyFile,
	}
	rbt := rabbitmq.NewRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("Unable to connect to rabbit", err)
	}

	// Consumer
	if app.IsConsumerMode() {
		cc := consumer.Config{
			ExchangeName:  mainConfig.Exchange.Name,
			ExchangeType:  mainConfig.Exchange.Type,
			RoutingKey:    mainConfig.Consumer.RoutingKey,
			QueueName:     mainConfig.Consumer.Queue,
			ConsumerName:  mainConfig.Consumer.ConnectionName,
			ConsumerCount: mainConfig.Consumer.Count,
			PrefetchCount: mainConfig.Consumer.PrefetchCount,
		}
		cc.Reconnect.MaxAttempt = 60
		cc.Reconnect.Interval = 1 * time.Second

		cc.Notifier = app.GetNotifier()

		csm := consumer.NewConsumer(cc, rbt)
		if err := csm.Start(); err != nil {
			log.Fatalln("Unable to start consumer!", err)
		}
		select {}

	}

	// Producer
	if app.IsProducerMode() {
		pc := producer.Config{
			ExchangeName: mainConfig.Exchange.Name,
			ExchangeType: mainConfig.Exchange.Type,
		}

		prd := producer.NewProducer(pc, rbt)
		if err := prd.Send(app.GetRoutingKey(), app.GetMessage()); err != nil {
			fmt.Println("Ошибка отправки сообщения", app.GetMessage(), err)
			os.Exit(1)
		}
	}

}
