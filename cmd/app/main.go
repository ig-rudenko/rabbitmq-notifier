package main

import (
	"log"
	"multiple-notifier/internal/config"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/internal/mode"
	"multiple-notifier/internal/producer"
	"multiple-notifier/internal/web"
	"multiple-notifier/pkg/rabbitmq"
	"os"
	"time"
)

func main() {
	mainConfig := config.NewConfig()
	app := mode.NewApp(mainConfig)
	app.ParseArgs()

	// RabbitMQ
	rc := rabbitmq.Config{
		Schema:     "amqps",
		Username:   app.Config.Rabbitmq.User,
		Password:   app.Config.Rabbitmq.Password,
		Host:       app.Config.Rabbitmq.Host,
		Port:       app.Config.Rabbitmq.Port,
		VHost:      app.Config.Rabbitmq.Vhost,
		CaCertFile: app.Config.Rabbitmq.CaCert,
		CertFile:   app.Config.Rabbitmq.CertFile,
		KeyFile:    app.Config.Rabbitmq.KeyFile,
	}
	rbt := rabbitmq.NewRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("Unable to connect to rabbit", err)
	}

	// Consumer
	if app.IsConsumerMode() {
		cc := consumer.Config{
			ExchangeName:  app.Config.Exchange.Name,
			ExchangeType:  app.Config.Exchange.Type,
			RoutingKey:    app.Config.Consumer.RoutingKey,
			QueueName:     app.Config.Consumer.Queue,
			ConsumerName:  app.Config.Consumer.ConnectionName,
			ConsumerCount: app.Config.Consumer.Count,
			PrefetchCount: app.Config.Consumer.PrefetchCount,
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
			ExchangeName: app.Config.Exchange.Name,
			ExchangeType: app.Config.Exchange.Type,
		}
		prd := producer.NewProducer(pc, rbt)
		if err := prd.CreateExchange(); err != nil {
			log.Fatalln("Не удалось создать exchange ->", err)
		}

		server := web.NewHttpServer(prd, app.Config.Producer.AuthToken)
		address := os.Args[2]
		server.Start(address)
	}

}
