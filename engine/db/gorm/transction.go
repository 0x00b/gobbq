package gorm

import (
	"context"

	"github.com/0x00b/gobbq/engine/db/trx"
)

type dbKeyType struct{}

var DBKey dbKeyType

type Transaction struct {
	db *GormDB
}

var _ trx.TransactionInf = &Transaction{}

func NewTransaction(db *GormDB) *Transaction {
	return &Transaction{db: db}
}

// for get GormDB
func (t *Transaction) DB(c context.Context) *GormDB {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		return db
	}
	return t.db
}

// transaction

func (t *Transaction) Begin(c context.Context) context.Context {
	db := t.db.begin()
	return context.WithValue(c, DBKey, db)
}

func (t *Transaction) Commit(c context.Context) {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		db.db.Commit()
	}
}

func (t *Transaction) Rollback(c context.Context, e error) {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		db.db.Rollback()
	}
}
