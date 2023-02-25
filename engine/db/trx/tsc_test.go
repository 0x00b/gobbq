package trx_test

import (
	"context"
	"testing"

	"github.com/0x00b/gobbq/engine/db/trx"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
)

const TTT = 1

// Test
//
//	@author jun
//	@date 2021-10-08 09:30:49
type Test struct{}

// Begin
//
//	@receiver Test
//	@param c
//	@return entity.Context
//	@author jun
//	@date 2021-10-08 09:27:56
func (Test) Begin(c entity.Context) entity.Context {
	// return context.WithValue(c, "test", "test")
	return nil
}

// Commit
//
//	@receiver Test
//	@param c
//	@author jun
//	@date 2021-10-08 09:30:56
func (Test) Commit(c entity.Context) {
	xlog.Traceln("Commit", c.Value("test"))

}

// Rollback
//
//	@receiver Test
//	@param c
//	@param e
//	@author jun
//	@date 2021-10-08 09:30:51
func (Test) Rollback(c entity.Context, e error) {
	xlog.Traceln("Rollback", c.Value("test"))
}

// test
//
//	@param c
//	@author jun
//	@date 2021-10-08 09:29:14
func test(c entity.Context) {

	ts := trx.Transaction{}
	ts.RegisterTransaction(Test{})
	_ = ts.Transaction(func(c entity.Context) error {
		xlog.Traceln("Transaction", c.Value("test"))
		return nil
	})(c)
	xlog.Traceln("test", c.Value("test"))
}

func TestTransaction(t *testing.T) {

	c := context.Background()

	// test(c)
	xlog.Traceln("TestTransaction", c.Value("test"))

}
