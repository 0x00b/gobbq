package trx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/db/trx"
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
//	@return context.Context
//	@author jun
//	@date 2021-10-08 09:27:56
func (Test) Begin(c context.Context) context.Context {
	return context.WithValue(c, "test", "test")

}

// Commit
//
//	@receiver Test
//	@param c
//	@author jun
//	@date 2021-10-08 09:30:56
func (Test) Commit(c context.Context) {
	fmt.Println("Commit", c.Value("test"))

}

// Rollback
//
//	@receiver Test
//	@param c
//	@param e
//	@author jun
//	@date 2021-10-08 09:30:51
func (Test) Rollback(c context.Context, e error) {
	fmt.Println("Rollback", c.Value("test"))
}

// test
//
//	@param c
//	@author jun
//	@date 2021-10-08 09:29:14
func test(c context.Context) {

	ts := trx.Transaction{}
	ts.RegisterTransaction(Test{})
	_ = ts.Transaction(func(c context.Context) error {
		fmt.Println("Transaction", c.Value("test"))
		return nil
	})(c)
	fmt.Println("test", c.Value("test"))
}

func TestTransaction(t *testing.T) {

	c := context.Background()

	test(c)
	fmt.Println("TestTransaction", c.Value("test"))

}
