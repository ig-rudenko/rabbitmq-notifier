package mode

import (
	"fmt"
	"log"
	"multiple-notifier/internal/config"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/internal/notifier/email"
	"multiple-notifier/internal/notifier/telegram"
	"os"
	"slices"
)

type App struct {
	Config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{config}
}

func (a *App) ParseArgs() {
	if len(os.Args) < 2 || (os.Args[1] != "consumer" && os.Args[1] != "producer") {
		a.ShowHelpText()
	}
	if os.Args[1] == "consumer" && len(os.Args) < 3 && !slices.Contains([]string{"telegram", "email"}, os.Args[2]) {
		a.ShowHelpText()
	}
	if os.Args[1] == "producer" && len(os.Args) < 3 {
		a.ShowHelpText()
	}
}

func (a *App) IsConsumerMode() bool {
	return os.Args[1] == "consumer"
}

func (a *App) IsProducerMode() bool {
	return os.Args[1] == "producer"
}

func (a *App) ShowHelpText() {
	fmt.Println("Для запуска необходимо передать параметр `consumer` или `producer`")
	fmt.Println("    `consumer` требует также следующий параметр - тип уведомителя. Доступны `telegram` `email` и `sms`.")
	log.Fatalln("    `producer` требует также следующий параметр - address, в формате 0.0.0.0:5555")
}

func (a *App) GetNotifier() consumer.Notifier {
	if os.Args[2] == "telegram" {
		return telegram.NewNotifier(a.Config.Consumer.ExpireAfterSeconds)
	}
	if os.Args[2] == "email" {
		return email.NewNotifier(a.Config.Consumer.ExpireAfterSeconds)
	}

	log.Fatalln("Неверный тип notifier")
	return nil
}
