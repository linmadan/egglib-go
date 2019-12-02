package sarama

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/core/application"
)

type OriginalMessageConverter struct {
}

func (converter *OriginalMessageConverter) ConvertToMessage(originalMessage interface{}) (*application.Message, error) {
	message := &application.Message{
	}
	if err := json.Unmarshal([]byte(originalMessage.(*sarama.ConsumerMessage).Value), &message); err != nil {
		return nil, err
	}
	return message, nil
}
