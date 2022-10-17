package rabbitmq

import (
	"encoding/json"
	"go-nsq/application/mq"
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/streadway/amqp"
)

type IRabbitMQClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RabbitMQClient struct {
	client *amqp.Connection
	msgCounter prometheus.Counter
}

func NewRabbitMQClient() (IRabbitMQClient, error) {
	url := os.Getenv("RABBITMQ_URL_ADDRESS")
	conn, err := amqp.Dial(url)
	if err != nil {

		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}
	reg := prometheus.NewRegistry()
	msgCounter := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name:      "RabbitMQ_message_pumped_count",
		Help:      "Number of message pumped by NSQ",
	})
	// Register msgCounter metric
	prometheus.Register(msgCounter)

	return RabbitMQClient{
		client: conn,
		msgCounter: msgCounter,
	}, nil
}

func (m RabbitMQClient) Publish(topic string, message []byte) error {
	ch, err := m.client.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		topic,
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

	//body := "Sending message from RabbitMQ"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {

		failOnError(err, "Failed to publish a message")
		return err
	}
	log.Printf(" [x] Congrats, sending message: %s", message)
	m.msgCounter.Inc()
	return nil
}

func (m RabbitMQClient) Subscribe(topic string) error {
	ch, err := m.client.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}

	q, err := ch.QueueDeclare(
		topic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}
	response, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}

	var data mq.Message
	for d := range response {
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Printf("Error unmarshalling json at RabbitMQ-Subscribe with error : %v", err)
			return err
		}
		log.Println("Logging message from RabbitMQ-Subscriber")
		log.Println(data.FileName)
		log.Println(data.FileObjectID)
		log.Println(data.Timestamp)
	}
	return nil
}
