package nsq

import (
	"log"

	nsq "github.com/nsqio/go-nsq"
)

type INSQClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

type Message struct {
	Timestamp    string
	FileObjectID string
	FileName     string
}

type NSQMessageHandler struct{}

func processMessage(body []byte) error {
	log.Println(body)
	return nil
}

func (h *NSQMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	err := processMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}

type NSQClient struct {
	config nsq.Config
}

func NewNSQClient() INSQClient {
	config := nsq.NewConfig()
	return &NSQClient{
		config: *config,
	}
}

func (n NSQClient) Publish(topic string, message []byte) error {
	//config := nsq.NewConfig()
	publisher, err := nsq.NewProducer("127.0.0.1:4150", &n.config)
	if err != nil {
		return err
	}

	err = publisher.Publish(topic, message)
	if err != nil {
		return err
	}

	return nil
}

func (n NSQClient) Subscribe(topic string) error {
	nsqSubscriber, err := nsq.NewConsumer(topic, "channel", &n.config)
	if err != nil {
		return err
	}
	nsqSubscriber.ConnectToNSQLookupd("localhost:4161")
	nsqSubscriber.AddHandler(&NSQMessageHandler{})

	return nil
}
