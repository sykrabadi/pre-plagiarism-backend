package nsq

import (
	nsq "github.com/nsqio/go-nsq"
)

type INSQClient interface {
	Publish(string, []byte) error
}

type Message struct {
	Timestamp    string
	FileObjectID string
	FileName     string
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
