package entity

import (
	"strconv"
	"sync/atomic"

	"github.com/0x00b/gobbq/engine/nets"
)

// EntityID proxyid + instid + id => (22bit + 10bit + 32bit)
type EntityID uint64

type ID uint64
type InstID uint32
type ProxyID uint32

const (
	TotalBitLen   = 64 // (proxy+inst+id)
	ProxyIDBitNum = 20
	InstIDBitNum  = 8
	IDBitNum      = 36

	ProxyIDMask = 1<<ProxyIDBitNum - 1
	InstIDMask  = 1<<InstIDBitNum - 1
	IDBitMask   = 1<<IDBitNum - 1
)

type EntityIDGenerator interface {
	NewEntityID() EntityID
}

func NewEntityID(pid ProxyID, iid InstID) EntityID {
	return FixedEntityID(pid, iid, GenID())
}

func FixedEntityID(pid ProxyID, iid InstID, id ID) EntityID {
	eid := uint64(pid)<<(TotalBitLen-ProxyIDBitNum) | uint64(iid&InstIDMask)<<IDBitNum | uint64(id&IDBitMask)
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
	return ProxyID((eid >> (TotalBitLen - ProxyIDBitNum)) & ProxyIDMask)
}

func (eid EntityID) InstID() InstID {
	return InstID((eid >> IDBitNum) & InstIDMask)
}

func (eid EntityID) ID() ID {
	return ID(eid & IDBitMask)
}

var u64IdCounter uint64

func GenID() ID {

	n := atomic.AddUint64(&u64IdCounter, 1) & IDBitMask

	return ID(n)
}
