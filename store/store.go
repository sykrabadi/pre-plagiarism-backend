package store

// store.go only contains interfaces

type Store interface {
	DocumentStore() DocumentStore
}

type DocumentStore interface {
	SendData() error
}
