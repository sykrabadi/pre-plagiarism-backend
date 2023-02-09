package mq

type MQPublishMessage struct {
	FileName     string      `json:"file_name"`
	FileObjectID interface{} `json:"file_object_id"`
	Timestamp    string      `json:"timestamp"`
}

type BoundingBox struct {
	X1 float64 `json:"x1"`
	X2 float64 `json:"x2"`
	Y1 float64 `json:"y1"`
	Y2 float64 `json:"y2"`
}

type MQSubscribeMessage struct {
	FileName      string        `json:"file_name"`
	FileObjectID  interface{}   `json:"file_object_id"`
	Timestamp     string        `json:"timestamp"`
	BoundingBoxes []BoundingBox `json:"bounding-box"`
}
