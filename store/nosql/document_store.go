package nosql

import (
	"context"
	"fmt"
	"go-nsq/application/mq"
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

func (c *DocumentStoreService) UpdateData(ctx context.Context, payload mq.MQSubscribeMessage) error {
	coll := c.conn.Db.Collection("documents")
	// TODO : Make contract to ensure the document schema
	id, err := primitive.ObjectIDFromHex(fmt.Sprint(payload.FileObjectID))
	if err != nil {
		log.Printf("[DocumentStoreService.UpdateData] error creating objectId with error % \n", err)
	}
	filter := bson.D{
		{"_id", id},
	}
	update := bson.D{
		{"$set", bson.D{
			{Key: "bounding_boxes", Value: bson.A{
				bson.D{
					{"rect", bson.D{
						{"x1", 123},
						{"x2", 123},
						{"y1", 123},
						{"y2", 123},
					}},
				},
				bson.D{
					{"rect", bson.D{
						{"x1", 123},
						{"x2", 123},
						{"y1", 123},
						{"y2", 123},
					}},
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
	if err != nil {
		log.Printf("[DocumentStoreService.UpdateData] unable to update data with object ID %v with error %v \n", payload.FileObjectID, err)
		return err
	}
	return nil
}
