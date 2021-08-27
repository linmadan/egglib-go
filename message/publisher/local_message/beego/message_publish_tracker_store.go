package beego

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/linmadan/egglib-go/message/publisher/local_message/beego/models"
	"github.com/linmadan/egglib-go/transaction/beego"
)

type MessagePublishTrackerStore struct {
	TransactionContext *beego.TransactionContext
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) getOrmer() orm.Ormer {
	var ormer orm.Ormer
	if messagePublishTrackerStore.TransactionContext != nil {
		ormer = messagePublishTrackerStore.TransactionContext.Ormer
	} else {
		ormer = orm.NewOrm()
	}
	return ormer
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) StartTrackMessagePublish() (int64, error) {
	ormer := messagePublishTrackerStore.getOrmer()
	messagePublishTrackerModel := &models.MessagePublishTracker{}
	if messagePublishTrackerId, err := ormer.Insert(messagePublishTrackerModel); err != nil {
		return 0, err
	} else {
		return messagePublishTrackerId, nil
	}
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) FinishTrackMessagePublish(messagePublishTrackerId int64) error {
	ormer := messagePublishTrackerStore.getOrmer()
	messagePublishTrackerModel := &models.MessagePublishTracker{
		Id: messagePublishTrackerId,
	}
	if _, err := ormer.Delete(messagePublishTrackerModel); err != nil {
		return err
	} else {
		return nil
	}
}

func (messagePublishTrackerStore *MessagePublishTrackerStore) HaveMessagePublishTracker() (bool, error) {
	ormer := messagePublishTrackerStore.getOrmer()
	querySeter := ormer.QueryTable("sys_message_publish_trackers")
	var messagePublishTrackerModels []*models.MessagePublishTracker
	if num, err := querySeter.All(&messagePublishTrackerModels); err != nil {
		return false, err
	} else {
		if num > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
}
