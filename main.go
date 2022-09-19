package main

import (
	"context"
	"go-nsq/db"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.DB.Disconnect(ctx)

	err = http.ListenAndServe(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
