package nosql

import (
	"go-nsq/db"
	"go-nsq/store"
)

type NoSQLStoreItem struct {
	documentStoreService DocumentStoreService
}

func NewNoSQLStore(db *db.Mongo) store.Store {
	return &NoSQLStoreItem{
		documentStoreService: DocumentStoreService{db},
	}
}

func (p *NoSQLStoreItem) DocumentStore() store.DocumentStore {
	return &p.documentStoreService
}
