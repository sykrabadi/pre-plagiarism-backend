package nosql

import (
	"context"
	"go-nsq/db"

	"go.mongodb.org/mongo-driver/bson"
)

type DocumentStoreService struct {
	conn *db.Mongo
}

func (c *DocumentStoreService) SendData() error {
	documentCollection := c.conn.Db.Collection("docs")
	_, err := documentCollection.InsertOne(context.Background(), bson.D{
		{Key: "name", Value: "TestInsertFromGo"},
	})

	if err != nil {
		return err
	}

	return nil
}
