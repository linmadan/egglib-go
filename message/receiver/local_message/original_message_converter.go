package local_message

import (
	"github.com/linmadan/egglib-go/core/application"
)

type OriginalMessageConverter interface {
	ConvertToMessage(originalMessage interface{}) (*application.Message, error)
}
