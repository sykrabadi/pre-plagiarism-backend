package store

import "context"

// store.go only contains interfaces

type Store interface {
	DocumentStore() DocumentStore
	MinioStore() MinioStore
}

type DocumentStore interface {
	SendData(string) error
	UpdateData(context.Context, string) error
}

type MinioStore interface {
	UploadDocument([]byte) error
}
