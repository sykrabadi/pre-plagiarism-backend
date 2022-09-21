package entrypoint

import (
	"context"
	"go-nsq/store"
)

type EntryPointService struct {
	DBStore store.Store
}

type IEntryPointService interface {
	SendData() error
	UpdateData(context.Context, string) error
}
