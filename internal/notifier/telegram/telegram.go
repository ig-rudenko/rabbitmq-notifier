package telegram

import (
	"bytes"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"multiple-notifier/internal/notifier"
	"net/http"
	"strconv"
)

type Notifier struct {
	*notifier.Notifier
}

type MessageData struct {
	ChatId    int    `json:"chatId"`
	Message   string `json:"message"`
	ParseMode string `json:"parseMode"`
	Token     string `json:"token"`
}

func NewNotifier(expireAfterSeconds uint64) *Notifier {
	return &Notifier{
		&notifier.Notifier{
			ExpireAfterSeconds: expireAfterSeconds,
			Name:               "TelegramNotifier",
		},
	}
}

func (n *Notifier) ProcessMessage(delivery *amqp.Delivery) bool {
	if n.IsDeliveryExpired(delivery) {
		return false
	}

	var messageData MessageData
	if !n.UnmarshallDelivery(delivery, &messageData) {
		return false
	}

	url := "https://api.telegram.org/bot" + messageData.Token + "/sendMessage"
	client := &http.Client{}
	values := map[string]string{"text": messageData.Message, "chat_id": strconv.Itoa(messageData.ChatId), "parse_mode": messageData.ParseMode}
	parameters, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(parameters))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		// Ошибка отправки телеграм сообщения.
		n.NegativeAcknowledgeDelivery(delivery)
		log.Printf("%s | Ошибка отправки телеграм сообщения. -> %e\n", n.Name, err)
		return false
	} else if res.StatusCode != 200 {
		// Неверный статус код отправки телеграм сообщения.
		n.NegativeAcknowledgeDelivery(delivery)
		log.Printf("%s | Проблема отправки телеграм сообщения. %s\n", n.Name, res.Status)
		defer res.Body.Close()
		return false
	} else {
		// Отправка успешна.
		log.Println(res.Status)
		n.AcknowledgeDelivery(delivery)
		defer res.Body.Close()
		return true
	}
}
