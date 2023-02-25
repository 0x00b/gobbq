package gorm_test

import (
	"testing"

	"github.com/0x00b/gobbq/engine/db/gorm"
	"github.com/0x00b/gobbq/xlog"
)

type B struct {
	*Cust
	Z int
}

type A struct {
	B
	X int
	Y int
}

type Cust struct {
	A
	X string
}

func TestT(t *testing.T) {

	cc := Cust{
		X: "ss",
		A: A{
			B: B{
				// Cust: &Cust{X: "zz"},
			},
		},
	}

	xlog.Traceln(gorm.ModelMap(nil, cc))

}
