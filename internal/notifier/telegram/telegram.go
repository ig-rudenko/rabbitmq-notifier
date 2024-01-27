package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"multiple-notifier/internal/misc"
	"net/http"
	"strconv"
	"time"
)

type Notifier struct {
	ExpireAfterSeconds uint64
}

type MessageData struct {
	ChatId    int    `json:"chatId"`
	Message   string `json:"message"`
	ParseMode string `json:"parseMode"`
	Token     string `json:"token"`
}

func NewNotifier(expireAfterSeconds uint64) *Notifier {
	return &Notifier{expireAfterSeconds}
}

func (n *Notifier) ProcessMessage(delivery *amqp.Delivery) bool {
	if n.ExpireAfterSeconds > 0 && delivery.Timestamp.Unix() < (time.Now().Unix()-int64(n.ExpireAfterSeconds)) {
		n.acknowledgeDelivery(delivery)
		return false
	}

	var messageData MessageData
	cleanedJson := misc.JsonEscape(string(delivery.Body))

	err := json.Unmarshal([]byte(cleanedJson), &messageData)
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
		defer res.Body.Close()
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
	if delivery.Redelivered {
		// Если сообщение уже второй раз обрабатывалось, то не даем ему бесконечно находиться в очереди.
		n.acknowledgeDelivery(delivery)
		return
	}
	err := delivery.Nack(false, true)
	if err != nil {
		fmt.Println("TelegramNotifier | Unable to requeue the message, dropped", err)
	}
}
