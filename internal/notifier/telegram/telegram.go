package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"strconv"
)

type Notifier struct {
}

type MessageData struct {
	ChatId    int    `json:"chatId"`
	Message   string `json:"message"`
	ParseMode string `json:"parseMode"`
	Token     string `json:"token"`
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) ProcessMessage(delivery *amqp.Delivery) bool {
	var messageData MessageData
	err := json.Unmarshal(delivery.Body, &messageData)
	if err != nil {
		// Неверный формат сообщения.
		fmt.Println("TelegramNotifier | Неверный формат сообщения ", err, delivery)
		// Закрываем сообщение, ведь его не получится уже обработать.
		n.acknowledgeDelivery(delivery)
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
		n.negativeAcknowledgeDelivery(delivery)
		fmt.Println("TelegramNotifier | Ошибка отправки телеграм сообщения.", err)
		return false
	} else if res.StatusCode != 200 {
		// Неверный статус код отправки телеграм сообщения.
		n.negativeAcknowledgeDelivery(delivery)
		fmt.Println("TelegramNotifier | Проблема отправки телеграм сообщения. StatusCode: ", res.Status)
		return false
	} else {
		// Отправка успешна.
		fmt.Println(res.Status)
		n.acknowledgeDelivery(delivery)
		defer res.Body.Close()
		return true
	}
}

func (n *Notifier) acknowledgeDelivery(delivery *amqp.Delivery) {
	err := delivery.Ack(false)
	if err != nil {
		fmt.Println("TelegramNotifier | Unable to acknowledge the message, dropped", err)
	}
}

func (n *Notifier) negativeAcknowledgeDelivery(delivery *amqp.Delivery) {
	err := delivery.Nack(false, true)
	if err != nil {
		fmt.Println("TelegramNotifier | Unable to requeue the message, dropped", err)
	}
}
