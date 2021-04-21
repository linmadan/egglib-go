package local_message

import (
	"fmt"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/message/publisher/local_message/beego"
	"github.com/linmadan/egglib-go/message/publisher/local_message/pg"
	beegoTransaction "github.com/linmadan/egglib-go/transaction/beego"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
)

type Publisher struct {
	messageStore MessageStore
}

func (publisher *Publisher) PublishMessages(messages []*application.Message, option map[string]interface{}) error {
	if dispatcher == nil {
		return fmt.Errorf("本地消息dispatcher还没有启动")
	}
	for _, message := range messages {
		if err := publisher.messageStore.AppendMessage(message); err != nil {
			return err
		}
	}
	if err := dispatcher.MessagePublishedNotice(); err != nil {
		return err
	}
	return nil
}

func NewLocalMessagePublisher(storeType string, storeOption map[string]interface{}) (*Publisher, error) {
	var messageStore MessageStore
	switch storeType {
	case "beego":
		var tc *beegoTransaction.TransactionContext
		if transactionContext, ok := storeOption["transactionContext"]; ok {
			tc = transactionContext.(*beegoTransaction.TransactionContext)
		} else {
			tc = nil
		}
		messageStore = &beego.MessagesStore{
			TransactionContext: tc,
		}
	case "pg":
		var tc *pgTransaction.TransactionContext
		if transactionContext, ok := storeOption["transactionContext"]; ok {
			tc = transactionContext.(*pgTransaction.TransactionContext)
		} else {
			tc = nil
		}
		messageStore = &pg.MessagesStore{
			TransactionContext: tc,
		}
	default:
		return nil, fmt.Errorf("无效的storeType: %s", storeType)
	}
	return &Publisher{
		messageStore: messageStore,
	}, nil
}
