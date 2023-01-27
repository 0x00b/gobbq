package gorm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/db/gorm"
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

	fmt.Println(gorm.ModelMap(context.Background(), cc))

}
