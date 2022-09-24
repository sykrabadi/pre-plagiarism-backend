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
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	location := "us-east-1"
	bucketName := os.Getenv("MINIO_BUCKET")
	err = client.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
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
	bucketName := os.Getenv("MINIO_BUCKET")

	info, err := m.client.PutObject(context.TODO(), bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(info.Location)
	log.Printf("Successfully uploaded at: %v \n", info.Location)
}
