package models

import "time"

type LocalMessage struct {
	tableName   string `pg:"sys_local_messages,alias:sys_local_message"`
	Id          int64  `pg:"pk:id"`
	MessageType string
	MessageBody string
	OccurredOn  time.Time
	IsPublished bool `pg:"default:false"`
}
