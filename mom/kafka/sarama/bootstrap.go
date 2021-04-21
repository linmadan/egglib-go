package sarama

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/linmadan/egglib-go/log"
	"strings"
	"sync"
)

func StartConsume(kafkaHosts string, groupId string, messageHandlerMap map[string]func(message *sarama.ConsumerMessage) error, logger log.Logger) error {
	config := sarama.NewConfig()
	//version, err := sarama.ParseKafkaVersion("2.1.1")
	//if err != nil {
	//	return err
	//}
	config.Version = sarama.V0_10_2_1
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	brokerList := strings.Split(kafkaHosts, ",")
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupId, config)
	if err != nil {
		return err
	}
	consumer := Consumer{
		ready:             make(chan bool),
		messageHandlerMap: messageHandlerMap,
		Logger:            logger,
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var topics []string
			for key := range messageHandlerMap {
				topics = append(topics, key)
			}
			if err := consumerGroup.Consume(ctx, topics, &consumer); err != nil {
				logger.Error(err.Error())
				return
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	logger.Info("sarama consumer启动成功，开始消费kafka消息")
	select {
	case <-ctx.Done():
		logger.Info("sarama consumer terminating, context cancelled")
	}
	cancel()
	wg.Wait()
	if err := consumerGroup.Close(); err != nil {
		return err
	}
	return nil
}
