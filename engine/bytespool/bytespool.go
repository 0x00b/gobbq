package bytespool

import (
	"fmt"
	"sync"
)

var (
	bufferPools = map[uint32]*sync.Pool{}
)

const (
	MaxBufferCap = 16 * 1024 * 1024 //16M
	MinBufferCap = 128              //128 byte

	bufferCapGrowMultiple = uint32(2) // twice
)

type Bytes struct {
	// must not resize the cap
	bytes []byte
}

func (bs *Bytes) Bytes() []byte {
	return bs.bytes
}

func init() {
	packetBodyCap := uint32(MinBufferCap)
	for packetBodyCap <= MaxBufferCap {
		key := CalcBufferCapKey(packetBodyCap)
		bufferPools[key] = &sync.Pool{
			New: func() interface{} {
				bs := Bytes{}
				bs.bytes = make([]byte, packetBodyCap)
				return bs
			},
		}
		packetBodyCap *= bufferCapGrowMultiple
	}
}

func CalcBufferCapKey(len uint32) uint32 {
	if len <= MinBufferCap {
		return 1
	}
	if len > MaxBufferCap {
		len = MaxBufferCap
	}
	var cnt uint32 = 1
	len -= 1
	len /= MinBufferCap
	for len != 0 {
		len /= 2
		cnt++
	}
	return cnt
}
func Get(len uint32) Bytes {
	key := CalcBufferCapKey(len)
	return bufferPools[key].Get().(Bytes)
}

func Put(bs Bytes) {
	capSize := uint32(cap(bs.bytes))
	if capSize > MaxBufferCap {
		fmt.Printf("bytes buffer len err: %d morethan %d\n", capSize, MaxBufferCap)
		capSize = MaxBufferCap
	}

	key := CalcBufferCapKey(capSize)
	bufferPools[key].Put(bs)
}
