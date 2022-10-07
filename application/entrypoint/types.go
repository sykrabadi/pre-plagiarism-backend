package entrypoint

import (
	"context"
	NSQ "go-nsq/application/mq/nsq"
	"go-nsq/application/mq/rabbitmq"
	"go-nsq/application/mq/redis"
	"go-nsq/store"
	"go-nsq/store/minio"
	"mime/multipart"
)

type EntryPointService struct {
	DBStore     store.Store
	NSQ         NSQ.INSQClient
	Minio       minio.MinioService
	RedisPubSub redis.IRedisClient
	RabbitMQ    rabbitmq.IRabbitMQClient
}

type IEntryPointService interface {
	SendData(*multipart.FileHeader) error
	UpdateData(context.Context, string) error
}
