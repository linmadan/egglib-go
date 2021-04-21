package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/message/publisher/local_message/pg/models"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
)

type MessagesStore struct {
	TransactionContext *pgTransaction.TransactionContext
}

func (messagesStore *MessagesStore) AppendMessage(message *application.Message) error {
	tx := messagesStore.TransactionContext.PgTx
	_, err := tx.QueryOne(
		pg.Scan(),
		"INSERT INTO sys_local_messages (id, message_type, message_body, occurred_on, is_published) VALUES (?, ?, ?, ?, ?)",
		message.MessageId, message.MessageType, message.MessageBody, message.OccurredOn, false)
	return err
}

func (messagesStore *MessagesStore) FindNoPublishedStoredMessages() ([]*application.Message, error) {
	tx := messagesStore.TransactionContext.PgTx
	var localMessageModels []*models.LocalMessage
	var messages []*application.Message
	query := tx.Model(&localMessageModels).Where("sys_local_message.is_published = ?", false)
	if err := query.Select(); err != nil {
		return nil, err
	}
	for _, localMessageModel := range localMessageModels {
		message := &application.Message{
			MessageId:   localMessageModel.Id,
			MessageType: localMessageModel.MessageType,
			MessageBody: localMessageModel.MessageBody,
			OccurredOn:  localMessageModel.OccurredOn,
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (messagesStore *MessagesStore) FinishPublishStoredMessages(messageIds []int64) error {
	tx := messagesStore.TransactionContext.PgTx
	_, err := tx.QueryOne(
		pg.Scan(),
		"UPDATE sys_local_messages SET is_published=? WHERE id IN (?)",
		true, pg.In(messageIds))
	return err
}
