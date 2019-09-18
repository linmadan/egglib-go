package domain

import "time"

type BaseEvent struct {
	OccurredOn time.Time `json:"occurredOn"`
}

type DomainEvent interface {
	EventType() string
}
