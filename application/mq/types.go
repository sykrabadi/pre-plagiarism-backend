package mq

import "time"

type Message struct {
	FileName     string
	FileObjectID string
	Timestamp    time.Duration
}
