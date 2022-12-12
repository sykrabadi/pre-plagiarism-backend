package entrypoint

import (
	"bytes"
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
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendToRest simulates data transfer using REST API
func SendToRest(data []byte) error{
	payload := bytes.NewBuffer(data)
	
	baseURL := "http://localhost:8082"
	
	resp, err := http.Post(baseURL+"/sendData", "application/json", payload)
	
	if err != nil {
		log.Fatalf("error sending data with error : %v", err)
	}
	defer resp.Body.Close()
	return nil
}

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
	// err = c.NSQ.Publish("TESTAGAIN", res)
	// if err != nil {
	// 	log.Printf("Error sending message to NSQ with error %v", err)
	// 	return err
	// }
	// Code below successfully
	// err = c.RedisPubSub.Publish("TESTAGAIN", res)
	// if err != nil {
	// 	log.Printf("Error sending message to RedisPubSub with error %v", err)
	// 	return err
	// }
	// err = c.RabbitMQ.Publish("TESTAGAIN", res)
	// if err != nil {
	// 	log.Printf("Error sending message to RabbitMQ with error %v", err)
	// 	return err
	// }
	// err = c.Kafka.Publish("TESTAGAIN", res)
	// if err != nil {
	// 	log.Printf("Error sending message to Kafka with error %v", err)
	// 	return err
	// }
	err = SendToRest(res)
	if err != nil {
		log.Printf("Error sending message to REST server with error %v", err)
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
