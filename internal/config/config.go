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

	// Значения по умолчанию.
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

	// Обновляем из переменных окружения.
	configuration.updateFromEnv()

	return &configuration
}

func (c *Config) updateFromEnv() {
	// Переопределяем через переменные окружения.
	c.Rabbitmq.User = misc.GetEnv("RABBITMQ_USER", c.Rabbitmq.User)
	c.Rabbitmq.Password = misc.GetEnv("RABBITMQ_PASSWORD", c.Rabbitmq.Password)
	c.Rabbitmq.Host = misc.GetEnv("RABBITMQ_HOST", c.Rabbitmq.Host)
	c.Rabbitmq.Port = misc.GetIntEnv("RABBITMQ_PORT", c.Rabbitmq.Port)
	c.Rabbitmq.Vhost = misc.GetEnv("RABBITMQ_VHOST", c.Rabbitmq.Vhost)
	c.Rabbitmq.CaCert = misc.GetEnv("RABBITMQ_CACERT", c.Rabbitmq.CaCert)
	c.Rabbitmq.CertFile = misc.GetEnv("RABBITMQ_CERTFILE", c.Rabbitmq.CertFile)
	c.Rabbitmq.KeyFile = misc.GetEnv("RABBITMQ_KEYFILE", c.Rabbitmq.KeyFile)

	c.Exchange.Name = misc.GetEnv("EXCHANGE_NAME", c.Exchange.Name)
	c.Exchange.Type = misc.GetEnv("EXCHANGE_TYPE", c.Exchange.Type)

	c.Consumer.ConnectionName = misc.GetEnv("CONSUMER_CONNECTION_NAME", c.Consumer.ConnectionName)
	c.Consumer.RoutingKey = misc.GetEnv("CONSUMER_ROUTING_KEY", c.Consumer.RoutingKey)
	c.Consumer.Queue = misc.GetEnv("CONSUMER_QUEUE", c.Consumer.Queue)
	c.Consumer.Count = misc.GetIntEnv("CONSUMER_COUNT", c.Consumer.Count)
	c.Consumer.PrefetchCount = misc.GetIntEnv("CONSUMER_PREFETCH_COUNT", c.Consumer.PrefetchCount)
	c.Consumer.ExpireAfterSeconds = misc.GetUIntEnv("CONSUMER_EXPIRE_AFTER_SECONDS", c.Consumer.ExpireAfterSeconds)

	c.Producer.AuthToken = misc.GetEnv("PRODUCER_AUTH_TOKEN", c.Producer.AuthToken)
}
