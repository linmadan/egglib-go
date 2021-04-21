package local_message

import (
	"fmt"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/message/receiver/local_message/beego"
	"github.com/linmadan/egglib-go/message/receiver/local_message/pg"
	"github.com/linmadan/egglib-go/message/receiver/local_message/sarama"
	beegoTransaction "github.com/linmadan/egglib-go/transaction/beego"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
)

type Receiver struct {
	OriginalMessageConverter OriginalMessageConverter
	ReceivedMessageStore     ReceivedMessageStore
}

func (receiver *Receiver) ReceiveMessage(originalMessage interface{}, option map[string]interface{}) (*application.Message, bool, error) {
	if message, err := receiver.OriginalMessageConverter.ConvertToMessage(originalMessage); err != nil {
		return nil, false, err
	} else {
		if receivedMessage, err := receiver.ReceivedMessageStore.FindMessage(message.MessageId); err != nil {
			return nil, false, err
		} else {
			if message.MessageId == receivedMessage.MessageId {
				return message, true, nil
			} else {
				return message, false, nil
			}
		}
	}
}

func (receiver *Receiver) ConfirmReceive(message *application.Message) error {
	if err := receiver.ReceivedMessageStore.SaveMessage(message); err != nil {
		return err
	}
	return nil
}

func NewLocalMessageReceiver(converterType string, converterOption map[string]interface{}, storeType string, storeOption map[string]interface{}) (*Receiver, error) {
	var originalMessageConverter OriginalMessageConverter
	switch converterType {
	case "sarama":
		originalMessageConverter = &sarama.OriginalMessageConverter{}
	default:
		return nil, fmt.Errorf("无效的originalMessageConverterType: %s", converterType)
	}
	var receivedMessageStore ReceivedMessageStore
	switch storeType {
	case "beego":
		var tc *beegoTransaction.TransactionContext
		if transactionContext, ok := storeOption["transactionContext"]; ok {
			tc = transactionContext.(*beegoTransaction.TransactionContext)
		} else {
			tc = nil
		}
		receivedMessageStore = &beego.ReceivedMessageStore{
			TransactionContext: tc,
		}
	case "pg":
		var tc *pgTransaction.TransactionContext
		if transactionContext, ok := storeOption["transactionContext"]; ok {
			tc = transactionContext.(*pgTransaction.TransactionContext)
		} else {
			tc = nil
		}
		receivedMessageStore = &pg.ReceivedMessageStore{
			TransactionContext: tc,
		}
	default:
		return nil, fmt.Errorf("无效的storeType: %s", storeType)
	}
	return &Receiver{
		OriginalMessageConverter: originalMessageConverter,
		ReceivedMessageStore:     receivedMessageStore,
	}, nil
}
