package local_message

import (
	"github.com/linmadan/egglib-go/core/application"
)

type MessageStore interface {
	AppendMessage(message *application.Message) error
	FindNoPublishedStoredMessages() ([]*application.Message, error)
	FinishPublishStoredMessages(messageIds []int64) error
}
