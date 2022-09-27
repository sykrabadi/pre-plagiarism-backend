package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type IRedisClient interface {
	Publish(string) error
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (IRedisClient, error) {
	config := redis.Options{
		Addr:     "localhost:6379",
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

func (r RedisClient) Publish(string) error {
	return nil
}
