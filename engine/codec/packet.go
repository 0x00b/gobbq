package codec

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/0x00b/gobbq/engine/bytespool"
	"github.com/0x00b/gobbq/proto/bbq"
)

// Packet is a packet for sending data
type Packet struct {
	refcount int32

	Src *PacketReadWriter // not nil indicates this is request packet

	// header: 只能在packet的生命周期内使用
	Header *bbq.Header

	totalLen  uint32
	headerLen uint32

	// data(header + body)
	bytes *bytespool.Bytes
}

const (
	minPacketBufferLen  = bytespool.MinBufferCap
	MaxPacketBodyLength = bytespool.MaxBufferCap
)

var (
	packetEndian = binary.LittleEndian
	packetPool   = &sync.Pool{
		New: func() interface{} {
			p := &Packet{
				Header: &bbq.Header{},
			}
			return p
		},
	}
)

type PacketType uint8

const (
	PacketRPC PacketType = 0x0
	PacketSys PacketType = 0x1
)

var packetName = map[PacketType]string{
	PacketRPC: "RPC",
	PacketSys: "PING",
}

func (t PacketType) String() string {
	if s, ok := packetName[t]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN_MESSAGE_TYPE_%d", uint8(t))
}

// Has reports whether f contains all (0 or more) flags in v.
func HasFlags(v, flags uint32) bool {
	return (v & flags) == flags
}

const (
	FlagDataChecksumIEEE uint32 = 0x01
)

// 获取packet的地方，作为函数返回值，强提醒记得release pkt
type releasePkt func()

func allocPacket() *Packet {
	pkt := packetPool.Get().(*Packet)

	pkt.reset()

	return pkt
}

// NewPacket allocates a new packet
func NewPacket() (*Packet, releasePkt) {
	pkt := allocPacket()
	return pkt, func() { pkt.Release() }
}

func (p *Packet) Retain() {
	atomic.AddInt32(&p.refcount, 1)
}

// Release releases the packet to packet pool
func (p *Packet) Release() {

	refcount := atomic.AddInt32(&p.refcount, -1)

	if refcount == 0 {
		bytespool.Put(p.bytes)
		packetPool.Put(p)
	} else if refcount < 0 {
		panic(fmt.Errorf("releasing packet with refcount=%d", p.refcount))
	}
}

// WriteBytes appends slice of bytes to the end of packetBody
func (p *Packet) WriteBody(b []byte) error {

	header, err := DefaultCodec.Marshal(p.Header)
	if err != nil {
		return err
	}

	p.headerLen = 4 + uint32(len(header))
	p.totalLen = p.headerLen + uint32(len(b))

	data := p.extendPacketData(p.totalLen)

	packetEndian.PutUint32(data, p.headerLen)
	copy(data[4:p.headerLen], header)
	copy(data[p.headerLen:p.totalLen], b)

	return nil
}

// PacketBody returns the total packetBody of packet
func (p *Packet) PacketBody() []byte {
	return p.bytes.Bytes()[p.headerLen:p.totalLen]
}

// PacketCap  returns the current packetBody capacity
func (p *Packet) GetPacketCap() uint32 {
	if p.bytes == nil {
		return 0
	}
	return uint32(cap(p.bytes.Bytes()))
}

func (p *Packet) reset() {
	p.Header.Reset()
	p.Src = nil
	p.headerLen = 0
	p.totalLen = 0
	p.refcount = 1
}

func (p *Packet) Data() []byte {
	if p.bytes == nil {
		return nil
	}
	return p.bytes.Bytes()[:p.totalLen]
}

// 返回的结果是在header的buf之后
func (p *Packet) extendPacketData(size uint32) []byte {
	if size > MaxPacketBodyLength {
		panic(ErrPacketBodyTooLarge)
	}

	oldCap := p.GetPacketCap()

	if size <= oldCap { // most case
		return p.Data()
	}

	// get new buffer

	bs := bytespool.Get(size)
	if bs == nil {
		panic("bytespool get bytes error")
	}

	// copy(bs.Bytes(), p.Data())
	oldBytes := p.bytes
	p.bytes = bs

	bytespool.Put(oldBytes)

	return p.Data()
}
