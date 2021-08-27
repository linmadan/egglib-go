package beego

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/linmadan/egglib-go/core/application"
	"github.com/linmadan/egglib-go/message/publisher/local_message/beego/models"
	"github.com/linmadan/egglib-go/transaction/beego"
	"strconv"
	"strings"
)

type MessagesStore struct {
	TransactionContext *beego.TransactionContext
}

func (messagesStore *MessagesStore) getOrmer() orm.Ormer {
	var ormer orm.Ormer
	if messagesStore.TransactionContext != nil {
		ormer = messagesStore.TransactionContext.Ormer
	} else {
		ormer = orm.NewOrm()
	}
	return ormer
}

func (messagesStore *MessagesStore) AppendMessage(message *application.Message) error {
	ormer := messagesStore.getOrmer()
	localMessageModel := &models.LocalMessage{
		Id:          message.MessageId,
		MessageType: message.MessageType,
		MessageBody: message.MessageBody,
		OccurredOn:  message.OccurredOn,
		IsPublished: false,
	}
	if _, err := ormer.Insert(localMessageModel); err != nil {
		return err
	}
	return nil
}

func (messagesStore *MessagesStore) FindNoPublishedStoredMessages() ([]*application.Message, error) {
	ormer := messagesStore.getOrmer()
	querySeter := ormer.QueryTable("sys_local_messages")
	var localMessageModels []*models.LocalMessage
	if _, err := querySeter.Filter("IsPublished", false).OrderBy("Id").All(&localMessageModels); err != nil {
		return nil, err
	}
	var messages []*application.Message
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
	var messageIdStrs []string
	for _, messageId := range messageIds {
		messageIdStrs = append(messageIdStrs, strconv.FormatInt(messageId, 10))
	}
	ormer := messagesStore.getOrmer()
	_, err := ormer.Raw("UPDATE sys_local_messages SET is_published = 1 WHERE Id IN (" + strings.Join(messageIdStrs, ",") + ")").Exec()
	if err != nil {
		return err
	}
	return nil
}
