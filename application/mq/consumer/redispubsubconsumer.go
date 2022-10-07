package consumer

import (
	redispubsub "go-nsq/application/mq/redis"
	"log"
)

func InitRedisPubSubSubscriber(
	client redispubsub.IRedisClient,
) {
	err := client.Subscribe("TESTAGAIN")
	if err != nil {
		log.Println(err)
		return
	}
}
