package models

import "time"

type MessagePublishTracker struct {
	tableName   string `pg:"sys_message_publish_trackers,alias:sys_message_publish_tracker"`
	Id        int64 `pg:"pk:id"`
	TrackTime time.Time `pg:"default:current_timestamp"`
}
