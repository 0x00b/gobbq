package entity

import (
	"math"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/0x00b/gobbq/engine/nets"
)

// EntityID proxyid + instid + id => (22bit + 10bit + 32bit)
type EntityID uint64

type ID uint32
type InstID uint32
type ProxyID uint32

const (
	proxyIDBitNum = 22
	instIDBitNum  = 10
	IdBitNum      = 32
	proxyIDMask   = 1<<proxyIDBitNum - 1
	instIDMask    = 1<<instIDBitNum - 1
	idBitMask     = 1<<IdBitNum - 1
)

type EntityIDGenerator interface {
	NewEntityID() EntityID
}

func NewEntityID(pid ProxyID, iid InstID) EntityID {
	return FixedEntityID(pid, iid, ID(GenIDU32()))
}

func FixedEntityID(pid ProxyID, iid InstID, id ID) EntityID {
	eid := uint64(pid)<<(64-proxyIDBitNum) | uint64(iid&instIDMask)<<IdBitNum | uint64(id&idBitMask)
	return EntityID(eid)
}

func DstEntity(pkt *nets.Packet) EntityID {
	return EntityID(pkt.Header.GetDstEntity())
}

func SrcEntity(pkt *nets.Packet) EntityID {
	return EntityID(pkt.Header.GetSrcEntity())
}

func (eid EntityID) Invalid() bool {
	return eid == 0
}

func (eid EntityID) String() string {
	return strconv.FormatUint(uint64(eid), 10)
}

func (eid EntityID) ProxyID() ProxyID {
	return ProxyID((eid >> (64 - proxyIDBitNum)) & proxyIDMask)
}

func (eid EntityID) InstID() InstID {
	return InstID((eid >> IdBitNum) & instIDMask)
}

func (eid EntityID) ID() ID {
	return ID(eid & idBitMask)
}

var u32IdCounter uint32

func GenIDU32() uint32 {

	i := atomic.AddUint32(&u32IdCounter, 1)
	n := uint32((time.Now().Unix()&math.MaxUint16)<<16) | uint32(i&math.MaxUint16)

	return uint32(n)
}
