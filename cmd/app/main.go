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

	fmt.Println(mainConfig)

	// RabbitMQ
	rc := rabbitmq.Config{
		Schema:         "amqps",
		Username:       mainConfig.RabbitmqUser,
		Password:       mainConfig.RabbitmqPass,
		Host:           mainConfig.RabbitmqHost,
		Port:           mainConfig.RabbitmqPort,
		VHost:          mainConfig.RabbitmqVhost,
		ConnectionName: mainConfig.RabbitmqConnectionName,
		CaCertFile:     mainConfig.CaCert,
		CertFile:       mainConfig.CertFile,
		KeyFile:        mainConfig.KeyFile,
	}
	rbt := rabbitmq.NewRabbit(rc)
	if err := rbt.Connect(); err != nil {
		log.Fatalln("Unable to connect to rabbit", err)
	}

	// Consumer
	if app.IsConsumerMode() {
		cc := consumer.Config{
			ExchangeName:  mainConfig.RabbitmqExchangeName,
			ExchangeType:  mainConfig.RabbitmqExchangeType,
			RoutingKey:    mainConfig.RabbitmqRoutingKey,
			QueueName:     mainConfig.RabbitmqQueue,
			ConsumerName:  mainConfig.RabbitmqConnectionName,
			ConsumerCount: mainConfig.ConsumerCount,
			PrefetchCount: mainConfig.PrefetchCount,
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
			ExchangeName: mainConfig.RabbitmqExchangeName,
			ExchangeType: mainConfig.RabbitmqExchangeType,
		}

		prd := producer.NewProducer(pc, rbt)
		if err := prd.Send(app.GetRoutingKey(), app.GetMessage()); err != nil {
			fmt.Println("Ошибка отправки сообщения", app.GetMessage(), err)
			os.Exit(1)
		}
	}

}
