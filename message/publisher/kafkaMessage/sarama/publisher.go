package sarama

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/core/application"
	"log"
	"strings"
	"time"
)

type Publisher struct {
	KafkaHosts string
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
			log.Fatalln(err)
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
				log.Printf("> message sent to topic %s partition %d at offset %d\n", message.MessageType, partition, offset)
			}
		}
	}
	return nil
}

func NewKafkaSaramaMessagePublisher(KafkaHosts string) (*Publisher, error) {
	return &Publisher{
		KafkaHosts: KafkaHosts,
	}, nil
}
