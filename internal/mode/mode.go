package mode

import (
	"fmt"
	"multiple-notifier/internal/consumer"
	"multiple-notifier/internal/notifier/telegram"
	"os"
	"slices"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) ParseArgs() {
	if len(os.Args) < 2 || (os.Args[1] != "consumer" && os.Args[1] != "producer") {
		a.ShowHelpText()
		os.Exit(1)
	}
	if os.Args[1] == "consumer" && len(os.Args) < 3 && !slices.Contains([]string{"telegram", "email", "sms"}, os.Args[2]) {
		a.ShowHelpText()
		os.Exit(1)
	}
	if os.Args[1] == "producer" && len(os.Args) < 4 {
		a.ShowHelpText()
		os.Exit(1)
	}
}

func (a *App) IsConsumerMode() bool {
	return os.Args[1] == "consumer"
}

func (a *App) IsProducerMode() bool {
	return os.Args[1] == "producer"
}

func (a *App) GetRoutingKey() string {
	return os.Args[2]
}

func (a *App) GetMessage() string {
	return os.Args[3]
}

func (a *App) ShowHelpText() {
	fmt.Println("Для запуска необходимо передать параметр `consumer` или `producer`")
	fmt.Println("    `consumer` требует также следующий параметр - тип уведомителя. Доступны `telegram` `email` и `sms`.")
	fmt.Println("    `producer` требует также следующие параметры - RoutingKey и JSON строку тела сообщения")
}

func (a *App) GetNotifier() consumer.Notifier {
	if os.Args[2] == "telegram" {
		return telegram.NewNotifier()
	}
	panic("Неверный тип notifier")
}