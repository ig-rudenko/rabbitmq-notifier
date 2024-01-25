package rabbitmq

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

type Config struct {
	Schema         string
	Username       string
	Password       string
	Host           string
	Port           string
	VHost          string
	ConnectionName string
}

type Rabbit struct {
	config     Config
	connection *amqp.Connection
}

// NewRabbit returns a RabbitMQ instance.
func NewRabbit(config Config) *Rabbit {
	return &Rabbit{
		config: config,
	}
}

// Connect connects to RabbitMQ server.
func (r *Rabbit) Connect() error {

	caCert, err := os.ReadFile("./rabbitmq-settings/root.crt")
	if err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(
		"./rabbitmq-settings/client.crt",
		"./rabbitmq-settings/client.key",
	)
	if err != nil {
		return err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)
	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
		ServerName:   "data.noc.sevtelecom.loc",
	}

	if r.connection == nil || r.connection.IsClosed() {
		con, err := amqp.DialTLS(fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s",
			r.config.Schema,
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
			r.config.VHost,
		), tlsConf)
		if err != nil {
			return err
		}
		r.connection = con
	}

	return nil
}

// Connection returns exiting `*amqp.Connection` instance.
func (r *Rabbit) Connection() (*amqp.Connection, error) {
	if r.connection == nil || r.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}

	return r.connection, nil
}

// Channel returns a new `*amqp.Channel` instance.
func (r *Rabbit) Channel() (*amqp.Channel, error) {
	chn, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}

	return chn, nil
}
