package main

import (
	"context"
	"go-nsq/application/entrypoint"
	"go-nsq/application/mq/consumer"
	"go-nsq/application/mq/kafka"
	nsqmq "go-nsq/application/mq/nsq"
	"go-nsq/application/mq/rabbitmq"
	"go-nsq/application/mq/redis"
	"go-nsq/db"
	"go-nsq/store/minio"
	"go-nsq/store/nosql"
	"go-nsq/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// TODO : Seperate current function calls to concurrent

func serveHTTP(
	serverAddr string,
	entrypointService entrypoint.IEntryPointService,
) {
	router := mux.NewRouter()
	serve := transport.NewHTTPServer(router, entrypointService)
	err := http.ListenAndServe(serverAddr, serve)
	if err != nil {
		log.Fatalf("Error connecting to %v", serverAddr)
	}
}

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
	redisPubSubClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Error intialize Redis Client")
	}
	
	// monitoringMetric := monitoring.InitMonitoring()
	rabbitMQClient, err := rabbitmq.NewRabbitMQClient()
	if err != nil {
		log.Fatalf("Error intialize RabbitMQ Client with error : %v", err)
	}
	kafkaClient, err := kafka.NewKafkaClient()
	if err != nil {
		log.Fatalf("Error intialize Kafka Client with error : %v", err)
	}
	entryPointService := entrypoint.NewEntryPointService(mongoDBStore, NSQClient, minio, redisPubSubClient, rabbitMQClient, kafkaClient)

	go func() {

		consumer.InitNSQSubscriber(NSQClient)
		consumer.InitRedisPubSubSubscriber(redisPubSubClient)
		consumer.InitRabbitMQSubscriber(rabbitMQClient)
	}()

	go serveHTTP(
		":8000",
		entryPointService,
	)
	
	consumer.InitKafkaSubscriber(kafkaClient)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("All server stopped!")
}
