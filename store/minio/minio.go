package minio

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client *minio.Client
	bucket string
}

type MinioService interface {
	UploadFile(string)
}

func InitMinioService(ctx context.Context, bucket string) (MinioService, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := true
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bucketLists, err := client.ListBuckets(context.TODO())
	log.Println(bucketLists)
	return &Minio{
		client: client,
		bucket: bucket,
	}, nil
}

func (m *Minio) UploadFile(objectName string) {
	file, err := os.Open(objectName)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}

	objectResult, err := m.client.FPutObject(context.TODO(), "documents", fileStat.Name(), `docker\minio\documents`, minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(objectResult.Location)
	log.Printf("Successfully uploaded at: %v \n", objectResult.Location)
}
