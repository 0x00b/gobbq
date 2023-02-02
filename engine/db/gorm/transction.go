package gorm

import (
	"github.com/0x00b/gobbq/engine/db/trx"
	"github.com/0x00b/gobbq/engine/entity"
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
func (t *Transaction) DB(c *entity.Context) *GormDB {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		return db
	}
	return t.db
}

// transaction

func (t *Transaction) Begin(c *entity.Context) *entity.Context {
	// db := t.db.begin()
	// return c. DBKey, db)
	return nil
}

func (t *Transaction) Commit(c *entity.Context) {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		db.db.Commit()
	}
}

func (t *Transaction) Rollback(c *entity.Context, e error) {
	db, ok := c.Value(DBKey).(*GormDB)
	if ok {
		db.db.Rollback()
	}
}
