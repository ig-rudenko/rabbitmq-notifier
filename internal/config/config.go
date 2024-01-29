package config

import (
	"encoding/json"
	"errors"
	"log"
	"multiple-notifier/internal/misc"
	"os"
)

type RabbitMQConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Vhost    string `json:"vhost,omitempty"`
	CaCert   string `json:"cacert"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

type ExchangeConfig struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

type ConsumerConfig struct {
	RoutingKey         string `json:"routingKey"`
	Queue              string `json:"queue"`
	ConnectionName     string `json:"connectionName"`
	Count              int    `json:"count,omitempty"`
	PrefetchCount      int    `json:"prefetchCount,omitempty"`
	ExpireAfterSeconds uint64 `json:"expireAfterSeconds,omitempty"`
}

type ProducerConfig struct {
	AuthToken string `json:"authToken,omitempty"`
}

type Config struct {
	Rabbitmq RabbitMQConfig `json:"rabbitmq"`
	Consumer ConsumerConfig `json:"consumer,omitempty"`
	Exchange ExchangeConfig `json:"exchange"`
	Producer ProducerConfig `json:"producer,omitempty"`
}

func NewConfig() *Config {
	configFilePath := misc.GetEnv("CONFIG_FILE", "/etc/rmq-notifier/config.json")
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Config file " + configFilePath + " does not exist")
	}

	file, _ := os.Open(configFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := json.NewDecoder(file)

	configuration := Config{
		Rabbitmq: RabbitMQConfig{
			Port:  5671,
			Vhost: "",
		},
		Exchange: ExchangeConfig{
			Type: "direct",
		},
		Consumer: ConsumerConfig{
			Count:         5,
			PrefetchCount: 5,
		},
		Producer: ProducerConfig{
			AuthToken: "",
		},
	}

	if err := decoder.Decode(&configuration); err != nil {
		log.Fatalln("INVALID CONFIG FILE", err)
	}

	// Переопределяем через переменную окружения.
	configuration.Producer.AuthToken = misc.GetEnv("PRODUCER_AUTH_TOKEN", configuration.Producer.AuthToken)

	return &configuration
}
