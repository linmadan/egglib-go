package application

import (
	"suplus-message/pkg/utils"
	"time"
)

type Message struct {
	MessageId   int64
	MessageType string
	MessageBody string
	OccurredOn  time.Time
}

func NewMessage(data map[string]interface{}) (*Message, error) {
	message := &Message{
		OccurredOn: time.Now(),
	}
	if IdWorker, err := utils.NewIdWorker(1); err != nil {
		return nil, err
	} else {
		if id, err := IdWorker.NextId(); err != nil {
			return nil, err
		} else {
			message.MessageId = id
		}
	}
	if messageType, ok := data["MessageType"]; ok {
		message.MessageType = messageType.(string)
	}
	if messageBody, ok := data["MessageBody"]; ok {
		message.MessageBody = messageBody.(string)
	}
	return message, nil
}
