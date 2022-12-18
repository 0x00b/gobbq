package codec

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/0x00b/gobbq/engine/bytespool"
)

// Packet is a packet for sending data
type Packet struct {
	refcount int32

	Src *PacketReadWriter // not nil indicates this is request packet

	// members:
	// _ uint32 // total len
	// _ PacketType
	// _ Flags
	// _ uint32 // request header len

	// data
	bytes *bytespool.Bytes
}

const (

	// packet = packet header + message header + message body
	// packet header = (4 byte total len) + (1 byte msg type) + (1 byte flag) + (4 byte header len)
	packetHeaderSize = 10

	minPacketBufferLen  = bytespool.MinBufferCap
	MaxPacketBodyLength = bytespool.MaxBufferCap - packetHeaderSize
)

var (
	packetEndian = binary.LittleEndian
	packetPool   = &sync.Pool{
		New: func() interface{} {
			p := &Packet{}
			p.bytes = bytespool.Get(minPacketBufferLen)
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

// depend on PacketType
type Flags uint8

// Has reports whether f contains all (0 or more) flags in v.
func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

// Has reports whether f contains all (0 or more) flags in v.
func (f *Flags) Set(v Flags) *Flags {
	*f = *f | v
	return f
}

const (
	//PacketRPC
	FlagDataRpcResponse  Flags = 0x01 // 0x00：request， 0x01：response
	FlagDataZipCompress  Flags = 0x02 // 留两位来代表压缩算法 {0x2,0x4,0x6}
	FlagDataChecksumIEEE Flags = 0x08
	FlagDataProtoBuf     Flags = 0x10
)

var flagName = map[PacketType]map[Flags]string{
	PacketRPC: {
		FlagDataRpcResponse:  "RpcResponse",
		FlagDataZipCompress:  "ZipCompress",
		FlagDataChecksumIEEE: "ChecksumIEEE",
		FlagDataProtoBuf:     "ProtoBuf",
	},
}

func allocPacket() *Packet {
	pkt := packetPool.Get().(*Packet)

	pkt.reset()

	return pkt
}

// NewPacket allocates a new packet
func NewPacket() *Packet {
	return allocPacket()
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
func (p *Packet) WriteBytes(b []byte) {
	pl := p.extendPacketBody(len(b))
	copy(pl, b)
}

// WriteBytes appends slice of bytes to the end of packetBody
func (p *Packet) WriteMsgHeader(b []byte) {
	p.setMsgHeaderLen(uint32(len(b)))
	p.WriteBytes(b)
}

// PacketBody returns the total packetBody of packet
func (p *Packet) PacketBody() []byte {
	return p.bytes.Bytes()[packetHeaderSize : packetHeaderSize+p.GetPacketBodyLen()]
}

// GetPacketBodyLen returns the packetBody length
func (p *Packet) GetPacketBodyLen() uint32 {
	// _ = p.bytes.Bytes()[3]
	return *(*uint32)(unsafe.Pointer(&p.bytes.Bytes()[0]))
}

func (p *Packet) setPacketBodyLen(plen uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes.Bytes()[0]))
	*pplen = plen
}

func (p *Packet) GetPacketType() PacketType {
	// _ = p.bytes.Bytes()[4]
	return *(*PacketType)(unsafe.Pointer(&p.bytes.Bytes()[4]))
}

func (p *Packet) SetPacketType(typ PacketType) {
	pplen := (*PacketType)(unsafe.Pointer(&p.bytes.Bytes()[4]))
	*pplen = typ
}

func (p *Packet) GetPacketFlags() Flags {
	// _ = p.bytes.Bytes()[5]
	return *(*Flags)(unsafe.Pointer(&p.bytes.Bytes()[5]))
}

func (p *Packet) SetPacketFlags(flags Flags) {
	pplen := (*Flags)(unsafe.Pointer(&p.bytes.Bytes()[5]))
	*pplen = flags
}

// GetPacketBodyLen returns the packetBody length
func (p *Packet) GetMsgHeaderLen() uint32 {
	// _ = p.bytes.Bytes()[3]
	return *(*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
}

func (p *Packet) setMsgHeaderLen(plen uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
	*pplen = plen
}

// PacketCap  returns the current packetBody capacity
func (p *Packet) PacketCap() uint32 {
	return uint32(len(p.bytes.Bytes()))
}

func (p *Packet) reset() {
	p.Src = nil
	p.setPacketBodyLen(0)
	p.SetPacketType(0)
	p.SetPacketFlags(0)
	p.refcount = 1
}

func (p *Packet) Data() []byte {
	return p.bytes.Bytes()[0 : packetHeaderSize+p.GetPacketBodyLen()]
}

func (p *Packet) packetBodySlice(i, j uint32) []byte {
	return p.bytes.Bytes()[i+packetHeaderSize : j+packetHeaderSize]
}

// 返回的结果是在header的buf之后
func (p *Packet) extendPacketBody(size int) []byte {
	if size > MaxPacketBodyLength {
		panic(ErrPacketBodyTooLarge)
	}

	packetBodyLen := p.GetPacketBodyLen()
	newPacketBodyLen := packetBodyLen + uint32(size)
	newPacketLen := packetHeaderSize + newPacketBodyLen
	oldCap := p.PacketCap()

	if newPacketLen <= oldCap { // most case
		p.setPacketBodyLen(newPacketBodyLen)
		return p.packetBodySlice(packetBodyLen, newPacketBodyLen)
	}

	// get new buffer

	if newPacketLen > MaxPacketBodyLength {
		panic(ErrPacketBodyTooLarge)
	}
	bs := bytespool.Get(newPacketLen)

	copy(bs.Bytes(), p.Data())
	oldBytes := p.bytes
	p.bytes = bs

	bytespool.Put(oldBytes)

	p.setPacketBodyLen(newPacketBodyLen)
	return p.packetBodySlice(packetBodyLen, newPacketBodyLen)
}
