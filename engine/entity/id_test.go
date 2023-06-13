package entity_test

import (
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/entity"
)

func TestXxx(t *testing.T) {

	id := entity.FixedEntityID(1111111111, 2222222222, 3333333333)

	fmt.Println(id.ProxyID(), id.InstID(), id.ID())

	id = 31243726709850116
	fmt.Println(id.ProxyID(), id.InstID(), id.ID())

}
