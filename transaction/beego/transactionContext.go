package beego

import "github.com/astaxie/beego/orm"

type TransactionContext struct {
	Ormer orm.Ormer
}

func (transactionContext *TransactionContext) StartTransaction() error {
	err := transactionContext.Ormer.Begin()
	return err
}

func (transactionContext *TransactionContext) CommitTransaction() error {
	err := transactionContext.Ormer.Commit()
	return err
}

func (transactionContext *TransactionContext) RollbackTransaction() error {
	err := transactionContext.Ormer.Rollback()
	return err
}

func NewBegoTransactionContext() *TransactionContext {
	return &TransactionContext{
		Ormer: orm.NewOrm(),
	}
}
