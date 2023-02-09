package store

import (
	"context"
	"go-nsq/application/mq"
)

// store.go only contains interfaces

type Store interface {
	DocumentStore() DocumentStore
}

type DocumentStore interface {
	SendData(string) (interface{}, error)
	UpdateData(context.Context, mq.MQSubscribeMessage) error
}
