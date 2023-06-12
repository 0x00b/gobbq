package entity

import (
	"math"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

// EntityID proxyid + instid + id => (16bit + 16bit + 32bit)
type EntityID uint64

type ID uint32
type InstID uint16
type ProxyID uint16

type EntityIDGenerator interface {
	NewEntityID() EntityID
}

func NewEntityID(pid ProxyID, iid InstID) EntityID {
	return FixedEntityID(pid, iid, ID(GenIDU32()))
}

func FixedEntityID(pid ProxyID, iid InstID, id ID) EntityID {
	eid := uint64(pid&math.MaxUint16)<<48 | uint64(iid&math.MaxUint16)<<32 | uint64(id&math.MaxUint32)
	return EntityID(eid)
}

func DstEntity(pkt *codec.Packet) EntityID {
	return EntityID(pkt.Header.GetDstEntity())
}

func SrcEntity(pkt *codec.Packet) EntityID {
	return EntityID(pkt.Header.GetSrcEntity())
}

func (eid EntityID) Invalid() bool {
	return eid == 0
}

func (eid EntityID) String() string {
	return strconv.FormatUint(uint64(eid), 10)
}

func (eid EntityID) ProxyID() ProxyID {
	return ProxyID(eid >> 48)
}

func (eid EntityID) InstID() InstID {
	return InstID((eid << 16) >> 48)
}

func (eid EntityID) ID() ID {
	return ID(eid & math.MaxUint32)
}

var u32IdCounter uint32

func GenIDU32() uint32 {

	i := atomic.AddUint32(&u32IdCounter, 1)
	n := uint32((time.Now().Unix()&math.MaxUint16)<<16) | uint32(i&math.MaxUint16)

	return uint32(n)
}
