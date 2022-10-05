package redis

import (
	"context"
	"time"

	"go-nsq/application/mq"

	"github.com/go-redis/redis/v9"
)

type Message struct {
	FileName     string
	FileObjectID string
	Timestamp    time.Duration
}

type IRedisClient interface {
	Publish(*mq.Message) error
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (IRedisClient, error) {
	config := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}

	rdb := redis.NewClient(&config)

	_, err := rdb.Ping(context.TODO()).Result()

	if err != nil {
		return nil, err
	}

	return RedisClient{
		client: rdb,
	}, nil
}

func (r RedisClient) Publish(Message *mq.Message) error {
	err := r.client.Publish(context.TODO(), "sendPDF", &Message).Err()

	if err != nil {
		return err
	}

	return nil
}
