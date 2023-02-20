package consumer

import (
	"go-nsq/application/mq/rabbitmq"
	"log"
)

func InitRabbitMQSubscriber(
	client rabbitmq.IRabbitMQClient,
) {
	err := client.Subscribe("update-document")
	if err != nil {
		log.Println(err)
		return
	}
}
