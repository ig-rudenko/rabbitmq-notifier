package notifier

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"multiple-notifier/internal/misc"
	"time"
)

type Notifier struct {
	ExpireAfterSeconds uint64
	Name               string
}

func (n *Notifier) IsDeliveryExpired(delivery *amqp.Delivery) bool {
	if n.ExpireAfterSeconds > 0 && delivery.Timestamp.Unix() < (time.Now().Unix()-int64(n.ExpireAfterSeconds)) {
		n.AcknowledgeDelivery(delivery)
		fmt.Printf("%s | Message has expired, dropped\n", n.Name)
		return true
	}
	return false
}

func (n *Notifier) UnmarshallDelivery(d *amqp.Delivery, v any) bool {
	cleanedJson := misc.JsonEscape(string(d.Body))
	err := json.Unmarshal([]byte(cleanedJson), v)
	if err != nil {
		// Неверный формат сообщения.
		fmt.Printf("%s | Неверный формат сообщения -> %e, %v\n", n.Name, err, d)
		// Закрываем сообщение, ведь его не получится уже обработать.
		n.AcknowledgeDelivery(d)
		return false
	}
	return true
}

func (n *Notifier) AcknowledgeDelivery(delivery *amqp.Delivery) {
	err := delivery.Ack(false)
	if err != nil {
		fmt.Printf("%s | Unable to acknowledge the message, dropped -> %e\n", n.Name, err)
	}
}

func (n *Notifier) NegativeAcknowledgeDelivery(delivery *amqp.Delivery) {
	if delivery.Redelivered {
		// Если сообщение уже второй раз обрабатывалось, то не даем ему бесконечно находиться в очереди.
		n.AcknowledgeDelivery(delivery)
		return
	}
	err := delivery.Nack(false, true)
	if err != nil {
		fmt.Printf("%s | Unable to requeue the message, dropped -> %e\n", n.Name, err)
	}
}
