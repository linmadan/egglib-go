package local_message

type MessagePublishTrackerStore interface {
	HaveMessagePublishTracker() (bool, error)
	StartTrackMessagePublish() (int64, error)
	FinishTrackMessagePublish(messagePublishTrackerId int64) error
}
