package consumer

import (
	"go-nsq/application/mq/kafka"
	"log"
)

func InitKafkaSubscriber(client kafka.IKafkaClient,){
	err := client.Subscribe("update-document")
	if err != nil{
		log.Printf("Error when subscribe to Kafka with topic : %v", err)
		return
	}
}
