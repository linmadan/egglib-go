package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/linmadan/egglib-go/message/publisher/local_message/pg/models"
	pgTransaction "github.com/linmadan/egglib-go/transaction/pg"
	"time"
)

type MessagePublishTrackerStore struct {
	TransactionContext *pgTransaction.TransactionContext
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) StartTrackMessagePublish() (int64, error) {
	tx := messagePublishTrackerStore.TransactionContext.PgTx
	var trackerId int64
	_, err := tx.QueryOne(
		pg.Scan(&trackerId),
		"INSERT INTO sys_message_publish_trackers (track_time) VALUES (?) RETURNING id",
		time.Now())
	return trackerId, err
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) FinishTrackMessagePublish(messagePublishTrackerId int64) error {
	tx := messagePublishTrackerStore.TransactionContext.PgTx
	messagePublishTrackerModel := new(models.MessagePublishTracker)
	messagePublishTrackerModel.Id = messagePublishTrackerId
	_, err := tx.Model(messagePublishTrackerModel).WherePK().Delete()
	return err
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) HaveMessagePublishTracker() (bool, error) {
	tx := messagePublishTrackerStore.TransactionContext.PgTx
	var messagePublishTrackerModels []*models.MessagePublishTracker
	query := tx.Model(&messagePublishTrackerModels)
	if count, err := query.SelectAndCount(); err != nil {
		return false, err
	} else {
		if count > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
}
