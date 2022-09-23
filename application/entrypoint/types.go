package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/db"
	"go-nsq/store"
)

type EntryPointService struct {
	DBStore store.Store
	MQ      mq.Client
	Minio   db.Minio
}

type IEntryPointService interface {
	SendData() error
	UpdateData(context.Context, string) error
}
