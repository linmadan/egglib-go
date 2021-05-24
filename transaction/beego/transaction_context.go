package beego

import "github.com/beego/beego/v2/client/orm"

type TransactionContext struct {
	Ormer   orm.Ormer
	TxOrmer orm.TxOrmer
}

func (transactionContext *TransactionContext) StartTransaction() error {
	var err error
	transactionContext.TxOrmer, err = transactionContext.Ormer.Begin()
	return err
}

func (transactionContext *TransactionContext) CommitTransaction() error {
	err := transactionContext.TxOrmer.Commit()
	return err
}

func (transactionContext *TransactionContext) RollbackTransaction() error {
	err := transactionContext.TxOrmer.Rollback()
	return err
}

func NewBeegoTransactionContext() *TransactionContext {
	return &TransactionContext{
		Ormer: orm.NewOrm(),
	}
}
