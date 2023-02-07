package model

type MQPublishMessage struct {
	Timestamp    string
	FileObjectID string
	FileName     string
}

type MQSubscribeMessage struct{}
