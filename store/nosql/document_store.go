package nosql

import (
	"context"
	"fmt"
	"go-nsq/db"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)
type DocumentStoreService struct {
	conn *db.Mongo
}

func (c *DocumentStoreService) SendData(documentName string) (string, error) {
	documentCollection := c.conn.Db.Collection("documents")
	res, err := documentCollection.InsertOne(context.Background(), bson.D{
		{Key: "name", Value: documentName},
	})
	fileObjectID := fmt.Sprint(res.InsertedID)

	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("Success insert document with ObjectID %v \n", res.InsertedID)
	return fileObjectID, nil
}

func (c *DocumentStoreService) UpdateData(ctx context.Context, objectID string) error {
	log.Println("[DocumentStoreService.UpdateData]")
	// documentCollection := c.conn.Db.Collection("docs")
	// TODO : Make contract to ensure the document schema
	// _, err := documentCollection.UpdateOne(ctx,
	// 	bson.D{},
	// 	bson.D{},
	// )

	// var result bson.M
	// findOptions := options.Find()
	// findOptions.SetSort(bson.D{{"_id", -1}})
	// findOptions.SetLimit(1)
	// cursor, err := documentCollection.Find(context.Background(), bson.D{}, findOptions)
	// if err != nil {
	// 	return err
	// }
	// cursor.Decode(&result)
	// output, err := json.MarshalIndent(result, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(output)
	return nil
}
