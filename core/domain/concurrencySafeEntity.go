package domain

type ConcurrencySafeEntity struct {
	Entity
	ConcurrencyVersion int `json:"concurrencyVersion"`
}
