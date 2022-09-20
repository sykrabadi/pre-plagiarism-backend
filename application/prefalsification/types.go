package prefalsification

import "go-nsq/store"

type PrefalsificationService struct {
	DBStore store.Store
}

type IPrefalsificationService interface {
	SendData() error
}
