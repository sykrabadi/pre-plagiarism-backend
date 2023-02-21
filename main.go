package main

import (
	"context"
	"fmt"
	"go-nsq/application/entrypoint"
	"go-nsq/application/mq/consumer"
	"go-nsq/application/mq/kafka"
	nsqmq "go-nsq/application/mq/nsq"
	"go-nsq/application/mq/rabbitmq"
	"go-nsq/application/mq/redis"
	"go-nsq/db"
	"go-nsq/externalapi/preplagiarism"
	"go-nsq/store/minio"
	"go-nsq/store/nosql"
	"go-nsq/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

// TODO : Seperate current function calls to concurrent

func serveHTTP(
	entrypointService entrypoint.IEntryPointService,
) {
	serverAddr := os.Getenv("SERVER_ADDR")
	router := mux.NewRouter()
	serve := transport.NewHTTPServer(router, entrypointService)
	err := http.ListenAndServe(fmt.Sprintf(":%v", serverAddr), serve)
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

	NSQClient := nsqmq.NewNSQClient(mongoDBStore)
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
	kafkaClient, err := kafka.NewKafkaClient(mongoDBStore)
	if err != nil {
		log.Fatalf("Error intialize Kafka Client with error : %v", err)
	}

	restyRestClient := resty.New()
	preplagiarismClient := preplagiarism.NewPrePlagiarismClient("http://localhost:8082", restyRestClient)
	if err != nil {
		log.Fatalf("error inizitialize REST Client with error : %v", err)
	}
	entryPointService := entrypoint.NewEntryPointService(mongoDBStore, NSQClient, minio, redisPubSubClient, rabbitMQClient, kafkaClient, preplagiarismClient)

	go func() {

		consumer.InitNSQSubscriber(NSQClient)
		consumer.InitRedisPubSubSubscriber(redisPubSubClient)
		consumer.InitRabbitMQSubscriber(rabbitMQClient)
	}()

	go serveHTTP(
		entryPointService,
	)
	
	consumer.InitKafkaSubscriber(kafkaClient)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("All server stopped!")
}
