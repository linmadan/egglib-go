package models

import "time"

type ReceivedMessage struct {
	tableName   string `pg:"sys_received_messages,alias:sys_received_message"`
	Id          int64  `pg:"pk:id"`
	MessageType string
	MessageBody string
	OccurredOn  time.Time
	ReceiveTime time.Time `pg:"default:current_timestamp"`
}
