package mq

type Message struct {
	FileName     string `json:"file_name"`
	FileObjectID string `json:"file_object_id"`
	Timestamp    string `json:"timestamp"`
}
