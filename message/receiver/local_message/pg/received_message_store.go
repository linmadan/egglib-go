package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/message/receiver/local_message/pg/models"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
	"time"
)

type ReceivedMessageStore struct {
	TransactionContext *pgTransaction.TransactionContext
}

func (receivedMessageStore *ReceivedMessageStore) SaveMessage(message *application.Message) error {
	tx := receivedMessageStore.TransactionContext.PgTx
	_, err := tx.QueryOne(
		pg.Scan(),
		"INSERT INTO sys_received_messages (id, message_type, message_body, occurred_on, receive_time) VALUES (?, ?, ?, ?, ?)",
		message.MessageId, message.MessageType, message.MessageBody, message.OccurredOn, time.Now())
	return err
}

func (receivedMessageStore *ReceivedMessageStore) FindMessage(messageId int64) (*application.Message, error) {
	tx := receivedMessageStore.TransactionContext.PgTx
	receivedMessageModel := new(models.ReceivedMessage)
	query := tx.Model(receivedMessageModel).Where("sys_received_message.id = ?", messageId)
	if err := query.First(); err != nil {
		if err.Error() != "pg: no rows in result set" {
			return nil, err
		}
	}
	message := &application.Message{
		MessageId:   receivedMessageModel.Id,
		MessageType: receivedMessageModel.MessageType,
		MessageBody: receivedMessageModel.MessageBody,
		OccurredOn:  receivedMessageModel.OccurredOn,
	}
	return message, nil
}
