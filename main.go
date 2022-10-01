package main

import (
	"context"
	"go-nsq/application/entrypoint"
	"go-nsq/application/mq/consumer"
	nsqmq "go-nsq/application/mq/nsq"
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
	redis, err := redis.NewRedisClient()
	entryPointService := entrypoint.NewEntryPointService(mongoDBStore, NSQClient, minio, redis)
	// server := transport.NewHTTPServer(entryPointService)
	// serverAddr := os.Getenv("SERVER_ADDR")
	// err = http.ListenAndServe(serverAddr, server)
	// if err != nil {
	// 	log.Fatalf("Error connect to the %s port \n", serverAddr)
	// }
	// consumer, err := nsq.NewConsumer("topic", "channel", nsq.NewConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Set the Handler for messages received by this Consumer. Can be called multiple times.
	// // See also AddConcurrentHandlers.
	// consumer.AddHandler(&nsqmq.NSQMessageHandler{})

	// // Use nsqlookupd to discover nsqd instances.
	// // See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	// err = consumer.ConnectToNSQLookupd("localhost:4161")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // wait for signal to exit
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan

	// // Gracefully stop the consumer.
	// consumer.Stop()

	consumer.InitNSQSubscriber(NSQClient)

	go serveHTTP(
		":8080",
		entryPointService,
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("All server stopped!")
}
