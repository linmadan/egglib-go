package local_message

import (
	"fmt"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/log"
	"github.com/linmadan/egglib-go/message/publisher/local_message/beego"
	"github.com/linmadan/egglib-go/message/publisher/local_message/pg"
	"github.com/linmadan/egglib-go/message/publisher/local_message/sarama"
	beegoTransaction "github.com/linmadan/egglib-go/transaction/beego"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
	"time"
)

var dispatcher *MessageDispatcher

type MessageDispatcher struct {
	notifications              chan struct{}
	dispatchTicker             *time.Ticker
	messageStore               MessageStore
	messagePublishTrackerStore MessagePublishTrackerStore
	messageEngine              MessageEngine
}

func (dispatcher *MessageDispatcher) MessagePublishedNotice() error {
	dispatcher.notifications <- struct{}{}
	return nil
}

func (dispatcher *MessageDispatcher) Dispatch() {
	for true {
		select {
		case <-dispatcher.dispatchTicker.C:
			go func(dispatcher *MessageDispatcher) {
				dispatcher.notifications <- struct{}{}
			}(dispatcher)
		case <-dispatcher.notifications:
			messages, messagePublishTrackerId, _ := dispatcher.loadNoPublishedMessages()
			if messagePublishTrackerId == -1 {
				go func(dispatcher *MessageDispatcher) {
					time.Sleep(time.Second * 1)
					dispatcher.notifications <- struct{}{}
				}(dispatcher)
			}
			if len(messages) > 0 {
				if messageIds, err := dispatcher.messageEngine.PublishToMessageSystem(messages, nil); err == nil {
					dispatcher.messageStore.FinishPublishStoredMessages(messageIds)
				}
			}
			if messagePublishTrackerId > 0 {
				dispatcher.messagePublishTrackerStore.FinishTrackMessagePublish(messagePublishTrackerId)
			}
		}
	}
}

func (dispatcher *MessageDispatcher) loadNoPublishedMessages() ([]*application.Message, int64, error) {
	if haveMessagePublishTracker, err := dispatcher.messagePublishTrackerStore.HaveMessagePublishTracker(); err != nil {
		return nil, 0, err
	} else {
		if haveMessagePublishTracker {
			return nil, -1, nil
		} else {
			messagePublishTrackerId, _ := dispatcher.messagePublishTrackerStore.StartTrackMessagePublish()
			messages, _ := dispatcher.messageStore.FindNoPublishedStoredMessages()
			return messages, messagePublishTrackerId, nil
		}
	}
}

func LaunchLocalMessageDispatcher(timeInterval time.Duration, messageEngineType string, engineOptions map[string]interface{}, storeType string, storeOptions map[string]interface{}, logger log.Logger) error {
	var messageEngine MessageEngine
	switch messageEngineType {
	case "sarama":
		var hosts string
		if kafkaHosts, ok := engineOptions["kafkaHosts"]; ok {
			hosts = kafkaHosts.(string)
		} else {
			hosts = "localhost:9092"
		}
		messageEngine = &sarama.Engine{
			KafkaHosts: hosts,
			Logger:     logger,
		}
	default:
		return fmt.Errorf("无效的messageEngineType: %s", messageEngineType)
	}
	var messageStore MessageStore
	var messagePublishTrackerStore MessagePublishTrackerStore
	switch storeType {
	case "beego":
		var tc *beegoTransaction.TransactionContext
		if transactionContext, ok := storeOptions["transactionContext"]; ok {
			tc = transactionContext.(*beegoTransaction.TransactionContext)
		} else {
			tc = nil
		}
		messageStore = &beego.MessagesStore{
			TransactionContext: tc,
		}
		messagePublishTrackerStore = &beego.MessagePublishTrackerStore{
			TransactionContext: tc,
		}
	case "pg":
		var tc *pgTransaction.TransactionContext
		if transactionContext, ok := storeOptions["transactionContext"]; ok {
			tc = transactionContext.(*pgTransaction.TransactionContext)
		} else {
			tc = nil
		}
		messageStore = &pg.MessagesStore{
			TransactionContext: tc,
		}
		messagePublishTrackerStore = &pg.MessagePublishTrackerStore{
			TransactionContext: tc,
		}
	default:
		return fmt.Errorf("无效的storeType: %s", storeType)
	}
	dispatcher = &MessageDispatcher{
		notifications:              make(chan struct{}),
		dispatchTicker:             time.NewTicker(timeInterval),
		messageStore:               messageStore,
		messagePublishTrackerStore: messagePublishTrackerStore,
		messageEngine:              messageEngine,
	}
	go func(dispatcher *MessageDispatcher) {
		dispatcher.Dispatch()
	}(dispatcher)
	return nil
}
