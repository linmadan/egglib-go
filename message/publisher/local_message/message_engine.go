package local_message

import (
	"github.com/linmadan/egglib-go/core/application"
)

type MessageEngine interface {
	PublishToMessageSystem(messages []*application.Message, option map[string]interface{}) ([]int64, error)
}
