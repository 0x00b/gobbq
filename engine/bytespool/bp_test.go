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
	packetBodyCap := uint32(bytespool.MinBufferCap)
	for packetBodyCap <= bytespool.MaxBufferCap {
		key := bytespool.CalcBufferCapKey(packetBodyCap)
		fmt.Println(i, key, packetBodyCap)
		i++
		packetBodyCap *= 2
	}
	fmt.Println("========")
	packetBodyCap = 0
	for packetBodyCap <= 1000 {
		key := bytespool.CalcBufferCapKey(packetBodyCap)
		fmt.Println(key, packetBodyCap)
		packetBodyCap++
	}

}
