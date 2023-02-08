package mq

type MQPublishMessage struct {
	FileName     string      `json:"file_name"`
	FileObjectID interface{} `json:"file_object_id"`
	Timestamp    string      `json:"timestamp"`
}

type MQSubscribeMessage struct {
	FileName     string      `json:"file_name"`
	FileObjectID interface{} `json:"file_object_id"`
	Timestamp    string      `json:"timestamp"`
	BoundingBox  map[string]float32
}
