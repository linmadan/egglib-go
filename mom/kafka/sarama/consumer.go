package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/log"
)

type Consumer struct {
	ready             chan bool
	messageHandlerMap map[string]func(message *sarama.ConsumerMessage) error
	Logger            log.Logger
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var append = make(map[string]interface{})
		append["framework"] = "sarama"
		append["topic"] = message.Topic
		append["messageTimestamp"] = message.Timestamp
		append["value"] = string(message.Value)
		consumer.Logger.Info("消费kafka消息", append)
		if err := consumer.messageHandlerMap[message.Topic](message); err == nil {
			session.MarkMessage(message, "")
		} else {
			var append = make(map[string]interface{})
			append["error"] = err.Error()
			consumer.Logger.Error("kafka消息处理错误", append)
		}
	}
	return nil
}
