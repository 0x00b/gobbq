package bytespool_test

import (
	"testing"

	"github.com/0x00b/gobbq/engine/bytespool"
	"github.com/0x00b/gobbq/xlog"
)

var ax []byte

func get() []byte {
	ax = make([]byte, 100)
	copy(ax, []byte("1211111111"))
	return ax
}
func test(bs []byte) {
	xlog.Traceln(len(bs))
	bs = bs[:50]
	xlog.Traceln(len(bs))
}

func TestMain(t *testing.T) {
	at := get()
	at = at[:10]
	copy(at, []byte("xxx"))
	test(ax)

	xlog.Traceln(len(ax), len(at), string(ax))

	//
	i := 0
	packetBodyCap := uint32(bytespool.MinBufferCap)
	for packetBodyCap <= bytespool.MaxBufferCap {
		key := bytespool.CalcBufferCapKey(packetBodyCap)
		xlog.Traceln(i, key, packetBodyCap)
		i++
		packetBodyCap *= 2
	}
	xlog.Traceln("========")
	packetBodyCap = 0
	for packetBodyCap <= 1000 {
		key := bytespool.CalcBufferCapKey(packetBodyCap)
		xlog.Traceln(key, packetBodyCap)
		packetBodyCap++
	}

}
