package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/store"
	"go-nsq/store/minio"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewEntryPointService(
	store store.Store,
	mq mq.Client,
	minio minio.MinioService,
) IEntryPointService {
	return &EntryPointService{
		DBStore: store,
		MQ:      mq,
		Minio:   minio,
	}
}

func (c *EntryPointService) SendData() error {
	docName := uuid.New().String()
	err := c.DBStore.DocumentStore().SendData(docName)

	if err != nil {
		return err
	}
	// config := nsq.NewConfig()
	// publisher, err := nsq.NewProducer("127.0.0.1:4150", config)
	// if err != nil {
	// 	return err
	// }
	c.Minio.UploadFile("2205.10133.pdf")

	msg := []byte("test publuish")
	err = c.MQ.Publish("TESTAGAIN", msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *EntryPointService) UpdateData(ctx context.Context, objectID string) error {
	fromHexID, _ := primitive.ObjectIDFromHex(objectID)
	id := primitive.ObjectID.String(fromHexID)
	err := c.DBStore.DocumentStore().UpdateData(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
