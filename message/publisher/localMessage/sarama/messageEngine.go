package sarama

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/core/application"
	"log"
	"strings"
	"time"
)

type Engine struct {
	KafkaHosts string
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
			log.Fatalln(err)
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
				log.Printf("FAILED to send message: %s\n", err)
			} else {
				log.Printf("> message sent to topic %s partition %d at offset %d\n", message.MessageType, partition, offset)
				successMessageIds = append(successMessageIds, message.MessageId)
			}
		}
	}
	return successMessageIds, nil
}
