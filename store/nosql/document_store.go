package nosql

import (
	"context"
	"fmt"
	"go-nsq/db"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type DocumentStoreService struct {
	conn *db.Mongo
}

func (c *DocumentStoreService) SendData(documentName string) (interface{}, error) {
	documentCollection := c.conn.Db.Collection("documents")
	res, err := documentCollection.InsertOne(context.Background(), bson.D{
		{Key: "name", Value: documentName},
	})
	//fileObjectID := fmt.Sprint(res.InsertedID)

	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("Success insert document with ObjectID %v \n", res.InsertedID)
	return res.InsertedID, nil
}

func (c *DocumentStoreService) UpdateData(ctx context.Context, objectID interface{}) error {
	coll := c.conn.Db.Collection("documents")
	// TODO : Make contract to ensure the document schema
	id, err := primitive.ObjectIDFromHex(fmt.Sprint(objectID))
	if err != nil{
		log.Printf("[DocumentStoreService.UpdateData] error creating objectId with error % \n", err)
	}
	filter := bson.D{
		{"_id", id},
	}
	update := bson.D{
		{"$set", bson.D{
			{"items", bson.A{
				"rect", bson.A{
					bson.D{
						{"x1", 123},
					},
					bson.D{
						{"y1", 321},
					},
					bson.D{
						{"x2", 1234},
					},
					bson.D{
						{"y2", 12345},
					},
				},
			},
		},
		},
	},
	}

	_, err = coll.UpdateOne(ctx,
		filter,
		update,
	)
	if err != nil{
		log.Printf("[DocumentStoreService.UpdateData] unable to update data with object ID %v with error %v \n", objectID, err)
		return err
	}

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
