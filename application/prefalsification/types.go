package prefalsification

import (
	"context"
	"go-nsq/store"
)

type PrefalsificationService struct {
	DBStore store.Store
}

type IPrefalsificationService interface {
	SendData() error
	UpdateData(context.Context, string) error
}
