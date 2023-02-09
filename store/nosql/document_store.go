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

	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("Success insert document with ObjectID %v \n", res.InsertedID)
	return res.InsertedID, nil
}

func (c *DocumentStoreService) UpdateData(ctx context.Context, payload mq.MQSubscribeMessage) error {
	coll := c.conn.Db.Collection("documents")
	id, err := primitive.ObjectIDFromHex(fmt.Sprint(payload.FileObjectID))
	if err != nil {
		log.Printf("[DocumentStoreService.UpdateData] error creating objectId with error % \n", err)
	}

	filter := bson.D{
		{"_id", id},
	}

	boundingBoxes := bson.A{}
	for idx, val := range payload.BoundingBoxes{
		log.Printf("[%v]:%v \n", idx, val)
		bb := bson.D{
			{"rect", bson.D{
				{"x1", val.X1},
				{"x2", val.X2},
				{"y1", val.Y1},
				{"y2", val.Y2},
			},
		},
		}
		boundingBoxes = append(boundingBoxes, bb)
	}

	update := bson.D{
		{"$set", bson.D{
			{Key: "bounding_boxes", Value: bson.A{
				boundingBoxes,
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
