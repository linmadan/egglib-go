package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(ReceivedMessage))
}

type ReceivedMessage struct {
	Id          int64 `orm:"pk"`
	MessageType string
	MessageBody string    `orm:"type(text)"`
	OccurredOn  time.Time `orm:"type(datetime)"`
	ReceiveTime time.Time `orm:"type(datetime)"`
}

func (receivedMessage *ReceivedMessage) TableName() string {
	return "sys_received_messages"
}
