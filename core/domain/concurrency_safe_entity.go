package domain

type ConcurrencySafeEntity interface {
	Entity
	ConcurrencyVersion() int64
}
