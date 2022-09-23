package db

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type Minio struct {
	Client *minio.Client
}

func loadMinioConfig() {}

func InitMinio(ctx context.Context) (*Minio, error) {
	client, err := minio.New("127.0.0.1:9000", &minio.Options{})
	if err != nil {
		return nil, err
	}
	return &Minio{
		Client: client,
	}, nil
}
