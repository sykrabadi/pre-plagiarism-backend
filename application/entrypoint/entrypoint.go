package entrypoint

import (
	"context"
	"encoding/json"
	"go-nsq/application/mq"
	"go-nsq/application/mq/kafka"
	nsqmq "go-nsq/application/mq/nsq"
	"go-nsq/application/mq/rabbitmq"
	"go-nsq/application/mq/redis"
	"go-nsq/store"
	"go-nsq/store/minio"
	"log"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewEntryPointService(
	store store.Store,
	nsq nsqmq.INSQClient,
	minio minio.MinioService,
	redisPubSub redis.IRedisClient,
	rabbitMQ rabbitmq.IRabbitMQClient,
	kafka kafka.IKafkaClient,
) IEntryPointService {
	return &EntryPointService{
		DBStore:     store,
		NSQ:         nsq,
		Minio:       minio,
		RedisPubSub: redisPubSub,
		RabbitMQ:    rabbitMQ,
		Kafka: kafka,
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
	message := mq.Message{
		Timestamp:    time.Now().String(),
		FileName:     fileName,
		FileObjectID: fileObjectID,
	}

	res, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling to json payload at EntryPointService-SendData")
		return err
	}
	err = c.NSQ.Publish("TESTAGAIN", res)
	if err != nil {
		log.Printf("Error sending message to NSQ with error %v", err)
		return err
	}
	// Code below successfully
	err = c.RedisPubSub.Publish("TESTAGAIN", res)
	if err != nil {
		log.Printf("Error sending message to RedisPubSub with error %v", err)
		return err
	}
	err = c.RabbitMQ.Publish("TESTAGAIN", res)
	if err != nil {
		log.Printf("Error sending message to RabbitMQ with error %v", err)
		return err
	}
	err = c.Kafka.Publish("TESTAGAIN", res)
	if err != nil {
		log.Printf("Error sending message to Kafka with error %v", err)
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
