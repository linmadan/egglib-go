package models

import (
	"github.com/astaxie/beego/orm"
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
	return "message_publish_trackers"
}
