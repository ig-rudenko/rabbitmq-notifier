package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

type Notifier struct {
	token  string
	chatId string
}

func NewNotifier(token, chatId string) *Notifier {
	return &Notifier{token, chatId}
}

func (n *Notifier) Send(message amqp.Delivery) {
	url := "https://api.telegram.org/bot" + n.token + "/sendMessage"

	client := &http.Client{}

	values := map[string]string{"text": string(message.Body), "chat_id": n.chatId}
	parameters, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(parameters))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
		defer res.Body.Close()
	}
}
