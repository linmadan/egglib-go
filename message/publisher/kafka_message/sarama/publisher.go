package sarama

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/log"
	"strings"
	"time"
)

type Publisher struct {
	KafkaHosts string
	Logger     log.Logger
}

func (publisher *Publisher) PublishMessages(messages []*application.Message, option map[string]interface{}) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokerList := strings.Split(publisher.KafkaHosts, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			publisher.Logger.Error(err.Error())
		}
	}()
	for _, message := range messages {
		if value, err := json.Marshal(message); err != nil {
			return err
		} else {
			msg := &sarama.ProducerMessage{
				Topic:     message.MessageType,
				Value:     sarama.StringEncoder(value),
				Timestamp: time.Now(),
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				return err
			} else {
				var append = make(map[string]interface{})
				append["framework"] = "sarama"
				append["topic"] = message.MessageType
				append["partition"] = partition
				append["offset"] = offset
				publisher.Logger.Info("生产kafka消息", append)
			}
		}
	}
	return nil
}

func NewKafkaSaramaMessagePublisher(kafkaHosts string, logger log.Logger) (*Publisher, error) {
	return &Publisher{
		KafkaHosts: kafkaHosts,
		Logger:     logger,
	}, nil
}
