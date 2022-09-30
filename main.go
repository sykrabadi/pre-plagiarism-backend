package main

import (
	"context"
	"go-nsq/application/entrypoint"
	nsqmq "go-nsq/application/mq/nsq"
	"go-nsq/application/mq/redis"
	"go-nsq/db"
	"go-nsq/store/minio"
	"go-nsq/store/nosql"
	"go-nsq/transport"
	"log"
	"net/http"
	"os"
)

// TODO : Seperate current function calls to concurrent

func main() {
	ctx := context.TODO()
	client, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.Client.Disconnect(ctx)

	mongoDBStore := nosql.NewNoSQLStore(client)

	NSQClient := nsqmq.NewNSQClient()
	minio, err := minio.InitMinioService(ctx, "documents")
	if err != nil {
		log.Fatalf("Error intialize Minio Client")
	}
	redis, err := redis.NewRedisClient()
	entryPointService := entrypoint.NewEntryPointService(mongoDBStore, NSQClient, minio, redis)
	server := transport.NewHTTPServer(entryPointService)
	serverAddr := os.Getenv("SERVER_ADDR")
	err = http.ListenAndServe(serverAddr, server)
	if err != nil {
		log.Fatalf("Error connect to the %s port \n", serverAddr)
	}
}
