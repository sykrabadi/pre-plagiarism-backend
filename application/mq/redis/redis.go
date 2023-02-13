package redis

import (
	"context"
	"encoding/json"
	"go-nsq/application/mq"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Message struct {
	FileName     string
	FileObjectID string
	Timestamp    time.Duration
}

type IRedisClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

type RedisClient struct {
	client *redis.Client
	msgCounter prometheus.Counter
}

func NewRedisClient() (IRedisClient, error) {
	config := redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	rdb := redis.NewClient(&config)

	_, err := rdb.Ping(context.TODO()).Result()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	reg := prometheus.NewRegistry()
	msgCounter := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name:      "RedisPubSub_message_pumped_count",
		Help:      "Number of message pumped by NSQ",
	})
	// Register msgCounter metric for redispubsub
	prometheus.Register(msgCounter)

	return RedisClient{
		client: rdb,
		msgCounter: msgCounter,
	}, nil
}

func (r RedisClient) Publish(Topic string, Message []byte) error {
	err := r.client.Publish(context.TODO(), "sendPDF", Message).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	r.msgCounter.Inc()
	return nil
}

func (r RedisClient) Subscribe(Channel string) error {
	subscriber := r.client.Subscribe(context.TODO(), Channel)

	// TODO : Fix subscription mechanism using subscriber.Channel to subscribe message concurrently
	msgs := subscriber.Channel()
	var resp mq.MQSubscribeMessage
	for d := range msgs {
		err := json.Unmarshal([]byte(d.Payload), &resp)
		if err != nil {
			log.Println("Fail to unmarshall json at Redis PubSub Subscription")
			return err
		}
		log.Printf("Logging message from Redis PubSub with payload : \n")
		log.Println(resp.FileName)
		log.Println(resp.FileObjectID)
		log.Println(resp.Timestamp)
	}

	return nil

}
