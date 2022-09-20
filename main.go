package main

import (
	"context"
	"go-nsq/application/prefalsification"
	"go-nsq/db"
	"go-nsq/store/nosql"
	"go-nsq/transport"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.Client.Disconnect(ctx)

	mongoDBStore := nosql.NewNoSQLStore(client)

	prefalsificationService := prefalsification.NewPrefalsificationService(mongoDBStore)
	server := transport.NewHTTPServer(prefalsificationService)
	serverAddr := os.Getenv("SERVER_ADDR")
	err = http.ListenAndServe(serverAddr, server)
	if err != nil {
		log.Fatalf("Error connect to the %s port \n", serverAddr)
	}
}
