package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/store"
	"go-nsq/store/minio"
)

type EntryPointService struct {
	DBStore store.Store
	MQ      mq.Client
	Minio   minio.MinioService
}

type IEntryPointService interface {
	SendData() error
	UpdateData(context.Context, string) error
}
