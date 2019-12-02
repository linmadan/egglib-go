package sarama

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/log"
	"strings"
	"time"
)

type Engine struct {
	KafkaHosts string
	Logger     log.Logger
}

func (engine *Engine) PublishToMessageSystem(messages []*application.Message, option map[string]interface{}) ([]int64, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokerList := strings.Split(engine.KafkaHosts, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			engine.Logger.Error(err.Error())
		}
	}()
	var successMessageIds []int64
	for _, message := range messages {
		if value, err := json.Marshal(message); err == nil {
			msg := &sarama.ProducerMessage{
				Topic:     message.MessageType,
				Value:     sarama.StringEncoder(value),
				Timestamp: time.Now(),
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				engine.Logger.Error(err.Error())
			} else {
				successMessageIds = append(successMessageIds, message.MessageId)
				var append = make(map[string]interface{})
				append["framework"] = "sarama"
				append["topic"] = message.MessageType
				append["partition"] = partition
				append["offset"] = offset
				engine.Logger.Info("kafka消息发送", append)
			}
		}
	}
	return successMessageIds, nil
}
