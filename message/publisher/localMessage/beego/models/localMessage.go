package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(LocalMessage))
}

type LocalMessage struct {
	Id          int64 `orm:"pk"`
	MessageType string
	MessageBody string    `orm:"type(text)"`
	OccurredOn  time.Time `orm:"type(datetime)"`
	IsPublished bool      `orm:"default(false)"`
}

func (localMessage *LocalMessage) TableName() string {
	return "sys_local_messages"
}
