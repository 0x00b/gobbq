package entity_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/tool/secure"
)

func TestXxx(t *testing.T) {

	c, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)
	secure.GO(func() {
		for {
			select {
			case <-c.Done():
				fmt.Println("done1")
				wg.Done()
				return
			}

		}

	})

	secure.GO(func() {
		for {
			select {
			case <-c.Done():
				fmt.Println("done2")
				wg.Done()
				return
			}

		}
	})

	cancel()
	wg.Wait()

	// id := entity.FixedEntityID(1111111111, 2222222222, 3333333333)

	// fmt.Println(id.ProxyID(), id.InstID(), id.ID())

	// id = 31243726709850116
	// fmt.Println(id.ProxyID(), id.InstID(), id.ID())

}
