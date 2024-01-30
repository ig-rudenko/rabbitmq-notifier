package email

import (
	"encoding/json"
	"errors"
	"log"
	"multiple-notifier/internal/misc"
	"os"
)

type NotifierConfig struct {
	Host     string `json:"host"`
	Port     uint64 `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type config struct {
	*NotifierConfig `json:"emailNotifier"`
}

func getConfig() *NotifierConfig {
	configFilePath := misc.GetEnv("CONFIG_FILE", "/etc/rmq-notifier/config.json")
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Config file " + configFilePath + " does not exist")
	}

	file, _ := os.Open(configFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	configuration := config{}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&configuration); err != nil {
		log.Fatalln("INVALID CONFIG FILE", err)
	}
	return configuration.NotifierConfig
}

func (c *config) updateFromEnv() {
	c.Host = misc.GetEnv("EMAIL_NOTIFIER_HOST", c.Host)
	c.Port = misc.GetUIntEnv("EMAIL_NOTIFIER_PORT", c.Port)
	c.Login = misc.GetEnv("EMAIL_NOTIFIER_LOGIN", c.Login)
	c.Password = misc.GetEnv("EMAIL_NOTIFIER_PASSWORD", c.Password)
}
