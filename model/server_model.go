package model

import "go.mongodb.org/mongo-driver/bson"

type ServerResponse struct {
	File_Name    *string `json:"file_name,omitempty"`
	ResponseCode int     `json:"responseCode,omitempty"`
	Message      string  `json:"Message,omitempty"`
	Data interface{} 
}

type BoundingBox struct {
	X1 float64 `json:"x1" bson:"x1"`
	X2 float64 `json:"x2" bson:"x2"`
	Y1 float64 `json:"y1" bson:"y1"`
	Y2 float64 `json:"y2" bson:"y2"`
}
type GetDocumentResponse struct {
	File_Name     string `json:"name" bson:"name"`
	BoundingBoxes []bson.A `json:"bounding_boxes" bson:"bounding_boxes"`
}
