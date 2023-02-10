package nosql

import (
	"context"
	"fmt"
	"go-nsq/application/mq"
	"go-nsq/db"
	"go-nsq/model"
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
	for _, val := range payload.BoundingBoxes{
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

func (c *DocumentStoreService) GetDocument(filename string) (*model.GetDocumentResponse, error){
	coll := c.conn.Db.Collection("documents")
	filter := bson.D{
		{Key:"name", Value:filename},
	}
	var res model.GetDocumentResponse
	err := coll.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil{
		log.Printf("[DocumentStoreService.GetDocument] error retrieving document from mongodb with error %v \n", err)
		return nil, err
	}
	log.Println(res)
	return &res, nil
}
