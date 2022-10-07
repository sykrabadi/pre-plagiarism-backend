package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type IRabbitMQClient interface {
	Publish() error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RabbitMQClient struct {
	client *amqp.Connection
}

func NewRabbitMQClient() (IRabbitMQClient, error) {
	url := os.Getenv("RABBITMQ_URL_ADDRESS")
	conn, err := amqp.Dial(url)
	if err != nil {

		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}

	return RabbitMQClient{
		client: conn,
	}, nil
}

func (m RabbitMQClient) Publish() error {
	ch, err := m.client.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"TESTAGAIN",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {

		failOnError(err, "Failed to declare a queue")
		return err
	}

	body := "Sending message from RabbitMQ"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {

		failOnError(err, "Failed to publish a message")
		return err
	}
	log.Printf(" [x] Congrats, sending message: %s", body)
	return nil
}
