package domain

import "errors"

type DomainEventPublisher interface {
	Reset() error
	Subscribe(domainEventSubscriber DomainEventSubscriber) error
	Publish(domainEvent DomainEvent) error
}

type BaseEventPublisher struct {
	isPublishing bool
	subscribers  []DomainEventSubscriber
}

func (domainEventPublisher *BaseEventPublisher) Reset() error {
	if domainEventPublisher.isPublishing {
		return errors.New("domain event is publishing, don't reset")
	}
	domainEventPublisher.subscribers = []DomainEventSubscriber{}
	return nil
}

func (domainEventPublisher *BaseEventPublisher) Subscribe(domainEventSubscriber DomainEventSubscriber) error {
	if domainEventPublisher.isPublishing {
		return errors.New("domain event is publishing, don't subscribe")
	}
	domainEventPublisher.subscribers = append(domainEventPublisher.subscribers, domainEventSubscriber)
	return nil
}

func (domainEventPublisher *BaseEventPublisher) Publish(domainEvent DomainEvent) error {
	if domainEventPublisher.isPublishing {
		return errors.New("domain event is publishing, don't publish")
	}
	domainEventPublisher.isPublishing = true
	for _, subscriber := range domainEventPublisher.subscribers {
		for _, eventType := range subscriber.SubscribedToEventTypes() {
			if eventType == domainEvent.EventType() {
				if err := subscriber.HandleEvent(domainEvent); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
