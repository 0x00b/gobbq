package bytespool_test

import (
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/bytespool"
)

var ax []byte

func get() []byte {
	ax = make([]byte, 100)
	copy(ax, []byte("1211111111"))
	return ax
}
func test(bs []byte) {
	fmt.Println(len(bs))
	bs = bs[:50]
	fmt.Println(len(bs))
}

func TestMain(t *testing.T) {
	at := get()
	at = at[:10]
	copy(at, []byte("xxx"))
	test(ax)

	fmt.Println(len(ax), len(at), string(ax))

	//
	i := 0
	messageBodyCap := uint32(bytespool.MinBufferCap)
	for messageBodyCap <= bytespool.MaxBufferCap {
		key := bytespool.CalcBufferCapKey(messageBodyCap)
		fmt.Println(i, key, messageBodyCap)
		i++
		messageBodyCap *= 2
	}
	fmt.Println("========")
	messageBodyCap = 0
	for messageBodyCap <= 1000 {
		key := bytespool.CalcBufferCapKey(messageBodyCap)
		fmt.Println(key, messageBodyCap)
		messageBodyCap++
	}

}
