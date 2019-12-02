package application

type (
	TransactionContext interface {
		StartTransaction() error
		CommitTransaction() error
		RollbackTransaction() error
	}
)
