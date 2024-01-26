package config

import (
	"encoding/json"
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
	RoutingKey     string `json:"routingKey"`
	Queue          string `json:"queue"`
	ConnectionName string `json:"connectionName"`
	Count          int    `json:"count,omitempty"`
	PrefetchCount  int    `json:"prefetchCount,omitempty"`
}

type Config struct {
	Rabbitmq RabbitMQConfig `json:"rabbitmq"`
	Consumer ConsumerConfig `json:"consumer"`
	Exchange ExchangeConfig `json:"exchange"`
}

func NewConfig() *Config {
	configFilePath := misc.GetEnv("CONFIG_FILE", "./config.json")
	file, _ := os.Open(configFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

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
	}

	if err := decoder.Decode(&configuration); err != nil {
		panic(err)
	}
	return &configuration
}
