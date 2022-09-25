package entrypoint

import (
	"context"
	"go-nsq/application/mq"
	"go-nsq/store"
	"go-nsq/store/minio"
	"mime/multipart"

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

func (c *EntryPointService) SendData(file *multipart.FileHeader) error {
	// TODO : Use fileObjectID as value to be sent to mq
	fileObjectID, err := c.DBStore.DocumentStore().SendData(file.Filename)

	if err != nil {
		return err
	}

	fileName, err := c.Minio.UploadFile(file)

	if err != nil {
		return err
	}

	msg := []byte(fileName)
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
