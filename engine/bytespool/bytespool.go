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
	bufferCap := uint32(MinBufferCap)
	for bufferCap <= MaxBufferCap {
		tempBufferCap := bufferCap
		key := CalcBufferCapKey(bufferCap)
		bufferPools[key] = &sync.Pool{
			New: func() interface{} {
				bs := &Bytes{}
				bs.bytes = make([]byte, tempBufferCap)
				return bs
			},
		}
		bufferCap *= bufferCapGrowMultiple
	}
}

func CalcBufferCapKey(len uint32) uint32 {
	if len == 0 {
		return 1
	}
	if len > MaxBufferCap {
		len = MaxBufferCap
	}
	var cnt uint32 = 1
	len -= 1
	len /= MinBufferCap
	for len != 0 {
		len /= bufferCapGrowMultiple
		cnt++
	}
	return cnt
}
func Get(len uint32) *Bytes {
	key := CalcBufferCapKey(len)
	return bufferPools[key].Get().(*Bytes)
}

func Put(bs *Bytes) {
	if bs == nil {
		return
	}
	capSize := uint32(cap(bs.bytes))
	if capSize > MaxBufferCap {
		panic(fmt.Sprintf("bytes buffer len err: %d morethan %d\n", capSize, MaxBufferCap))
	}

	key := CalcBufferCapKey(capSize)
	bufferPools[key].Put(bs)
}
