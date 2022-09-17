package main

import (
	"context"
	"go-nsq/db"
	"log"
)

func main() {
	ctx := context.Background()
	client, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
