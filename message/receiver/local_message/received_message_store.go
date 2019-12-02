package local_message

import (
	"github.com/linmadan/egglib-go/core/application"
)

type ReceivedMessageStore interface {
	SaveMessage(message *application.Message) error
	FindMessage(messageId int64) (*application.Message, error)
}
