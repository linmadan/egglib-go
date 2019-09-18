package domain

type DomainEventSubscriber interface {
	HandleEvent(domainEvent DomainEvent) error
	SubscribedToEventTypes() []string
}