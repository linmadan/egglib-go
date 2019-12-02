package local_message

import (
	"time"
)

type MessagePublishTracker struct {
	TrackerId int64
	TrackTime time.Time
}
