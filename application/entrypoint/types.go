package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/store"
)

type EntryPointService struct {
	DBStore store.Store
	MQ      mq.Client
}

type IEntryPointService interface {
	SendData() error
	UpdateData(context.Context, string) error
}
