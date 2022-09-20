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

func (c *DocumentStoreService) UpdateData(ctx context.Context, objectID string) error {
	documentCollection := c.conn.Db.Collection("docs")
	// TODO : Make contract to ensure the document schema
	_, err := documentCollection.UpdateOne(ctx,
		bson.D{},
		bson.D{},
	)

	if err != nil {
		return err
	}
	return nil
}
