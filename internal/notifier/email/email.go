package email

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
	"multiple-notifier/internal/notifier"
	"net/mail"
)

type Mail struct {
	Sender  string   `json:"sender"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Notifier struct {
	*notifier.Notifier
	Config *NotifierConfig
}

func NewNotifier(expireAfterSeconds uint64) *Notifier {
	return &Notifier{
		&notifier.Notifier{
			ExpireAfterSeconds: expireAfterSeconds,
			Name:               "EmailNotifier",
		},
		getConfig(),
	}
}

func (n *Notifier) ProcessMessage(delivery *amqp.Delivery) bool {
	if n.IsDeliveryExpired(delivery) {
		return false
	}

	mailData := Mail{}
	if !n.UnmarshallDelivery(delivery, &mailData) {
		return false
	}

	if err := n.sendEmail(mailData); err != nil {
		n.NegativeAcknowledgeDelivery(delivery)
		return false
	}

	return true
}

func (n *Notifier) sendEmail(mailData Mail) error {
	from := mail.Address{Address: n.Config.Login}
	to := make([]string, len(mailData.To))
	for i, addr := range mailData.To {
		a := mail.Address{Address: addr}
		to[i] = a.String()
	}

	message := gomail.NewMessage()
	message.SetHeader("From", from.String())
	message.SetHeader("To", to...)
	message.SetHeader("Subject", mailData.Subject)
	message.SetBody("text/html", mailData.Body)

	d := gomail.NewDialer(n.Config.Host, int(n.Config.Port), n.Config.Login, n.Config.Password)

	err := d.DialAndSend(message)
	return err
}
