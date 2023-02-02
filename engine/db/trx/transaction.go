package trx

import "github.com/0x00b/gobbq/engine/entity"

type Transaction []TransactionInf

type TransactionInf interface {
	Begin(c *entity.Context) *entity.Context
	Commit(c *entity.Context)
	Rollback(c *entity.Context, e error)
}

// Transaction 开启一个事务 f().
// useage: e := ts.Transaction(&c, func()error{return nil})
// 传进来的c会被修改，如果希望事务结束之后不影响c，那么需要：
// ctx := c
// e := ts.Transaction(&ctx, func()error{return nil})
func (ts *Transaction) Transaction(f func(*entity.Context) error) func(*entity.Context) error {
	return func(c *entity.Context) (e error) {
		for _, t := range *ts {
			c = t.Begin(c)
		}
		defer func() {
			if e != nil {
				for _, t := range *ts {
					t.Rollback(c, e)
				}
			} else {
				for _, t := range *ts {
					t.Commit(c)
				}
			}
		}()
		return f(c)
	}
}

func (ts *Transaction) RegisterTransaction(t TransactionInf) {
	*ts = append(*ts, t)
}
