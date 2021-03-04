package pg

import (
	"github.com/go-pg/pg/v10"
)

type TransactionContext struct {
	PgDd *pg.DB
	PgTx *pg.Tx
}

func (transactionContext *TransactionContext) StartTransaction() error {
	tx, err := transactionContext.PgDd.Begin()
	if err != nil {
		return err
	}
	transactionContext.PgTx = tx
	return nil
}

func (transactionContext *TransactionContext) CommitTransaction() error {
	err := transactionContext.PgTx.Commit()
	return err
}

func (transactionContext *TransactionContext) RollbackTransaction() error {
	err := transactionContext.PgTx.Rollback()
	return err
}

func NewPGTransactionContext(pgDd *pg.DB) *TransactionContext {
	return &TransactionContext{
		PgDd: pgDd,
	}
}
