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
	Port     uint16 `json:"port"`
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
