package consumer

import (
	redispubsub "go-nsq/application/mq/redis"
	"log"
)

func InitRedisPubSubSubscriber(
	client redispubsub.IRedisClient,
) {
	go func() {
		err := client.Subscribe("sendPDF")
		if err != nil {
			log.Println(err)
			return
		}
	}()
}
