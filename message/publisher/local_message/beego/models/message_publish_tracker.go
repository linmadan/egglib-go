package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(MessagePublishTracker))
}

type MessagePublishTracker struct {
	Id        int64
	TrackTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func (publishedMessageTracker *MessagePublishTracker) TableName() string {
	return "sys_message_publish_trackers"
}
