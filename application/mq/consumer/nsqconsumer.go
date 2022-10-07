package consumer

import (
	nsqmq "go-nsq/application/mq/nsq"
	"log"
)

type Handler struct {
}

func InitNSQSubscriber(
	client nsqmq.INSQClient,
) {
	err := client.Subscribe("TESTAGAIN")
	if err != nil {
		log.Println(err)
		return
	}
}
