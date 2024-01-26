package config

import (
	"encoding/json"
	"multiple-notifier/internal/misc"
	"os"
)

type Config struct {
	CaCert                 string `json:"cacert"`
	CertFile               string `json:"certFile"`
	KeyFile                string `json:"keyFile"`
	RabbitmqUser           string `json:"rabbitmqUser"`
	RabbitmqPass           string `json:"rabbitmqPass"`
	RabbitmqHost           string `json:"rabbitmqHost"`
	RabbitmqPort           int    `json:"rabbitmqPort,omitempty"`
	RabbitmqVhost          string `json:"rabbitmqVhost,omitempty"`
	RabbitmqConnectionName string `json:"rabbitmqConnectionName"`
	RabbitmqExchangeName   string `json:"rabbitmqExchangeName"`
	RabbitmqExchangeType   string `json:"rabbitmqExchangeType,omitempty"`
	RabbitmqRoutingKey     string `json:"rabbitmqRoutingKey"`
	RabbitmqQueue          string `json:"rabbitmqQueue"`
	ConsumerCount          int    `json:"consumerCount,omitempty"`
	PrefetchCount          int    `json:"prefetchCount,omitempty"`
}

func NewConfig() *Config {
	configFilePath := misc.GetEnv("CONFIG_FILE", "./config.json")
	file, _ := os.Open(configFilePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	configuration := Config{
		RabbitmqPort:         5671,
		RabbitmqVhost:        "",
		RabbitmqExchangeType: "direct",
		ConsumerCount:        5,
		PrefetchCount:        5,
	}

	if err := decoder.Decode(&configuration); err != nil {
		panic(err)
	}
	return &configuration
}
