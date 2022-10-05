package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Message struct {
	FileName     string
	FileObjectID string
	Timestamp    time.Duration
}

type IRedisClient interface {
	Publish(string, []byte) error
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (IRedisClient, error) {
	config := redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	rdb := redis.NewClient(&config)

	_, err := rdb.Ping(context.TODO()).Result()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return RedisClient{
		client: rdb,
	}, nil
}

func (r RedisClient) Publish(Topic string, Message []byte) error {
	err := r.client.Publish(context.TODO(), "sendPDF", Message).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
