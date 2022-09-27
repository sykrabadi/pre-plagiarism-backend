package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/application/mq/redis"
	"go-nsq/store"
	"go-nsq/store/minio"
	"mime/multipart"
)

type EntryPointService struct {
	DBStore     store.Store
	MQ          mq.Client
	Minio       minio.MinioService
	RedisPubSub redis.IRedisClient
}

type IEntryPointService interface {
	SendData(*multipart.FileHeader) error
	UpdateData(context.Context, string) error
}
