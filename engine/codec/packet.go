package codec

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"

	"github.com/0x00b/gobbq/engine/bytespool"
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

const (
	// (4 byte total len) + (1 byte msg type) + (1 byte flag) + (4 byte msg id)
	packetHeaderSize    = 10
	minPacketBufferLen  = bytespool.MinBufferCap
	MaxPacketBodyLength = bytespool.MaxBufferCap - packetHeaderSize
)

type PacketType uint8

const (
	PacketRPC  PacketType = 0x0
	PacketPing PacketType = 0x1
)

var packetName = map[PacketType]string{
	PacketRPC:  "RPC",
	PacketPing: "PING",
}

func (t PacketType) String() string {
	if s, ok := packetName[t]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN_FRAME_TYPE_%d", uint8(t))
}

type packetParser func(pc *PacketReadWriter, pkt *Packet) (uint32, error)

var packetParsers map[PacketType]packetParser

// depend on PacketType
type Flags uint8

const (
	//PacketRPC
	FlagDataZipCompress  Flags = 0x1
	FlagDataChecksumIEEE Flags = 0x2
	FlagDataProtoBuf     Flags = 0x4 //use proto
)

var flagName = map[PacketType]map[Flags]string{
	PacketRPC: {
		FlagDataZipCompress:  "ZipCompress",
		FlagDataChecksumIEEE: "ChecksumIEEE",
		FlagDataProtoBuf:     "ProtoBuf",
	},
}

func init() {
	packetParsers = map[PacketType]packetParser{
		PacketRPC: rpcPacketParser,
	}
}

func rpcPacketParser(pc *PacketReadWriter, pkt *Packet) (uint32, error) {
	// var buff16uint [2]byte
	// io.ReadFull(pc.rw, buff16uint[:])
	// mothodLen := packetEndian.Uint16(buff16uint[:])
	// buf := make([]byte, mothodLen)
	// _, err := io.ReadFull(pc.rw, buf)
	// if err != nil {
	// 	return 0, err
	// }
	// pkt.method = string(buf)
	// return 2 + uint32(mothodLen), nil
	return 0, nil
}

// Has reports whether f contains all (0 or more) flags in v.
func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

// Packet is a packet for sending data
type Packet struct {
	src *PacketReadWriter

	// depend on PacketType
	meta interface{}

	bytes bytespool.Bytes
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

// WriteBytes appends slice of bytes to the end of packetBody
func (p *Packet) WriteBytes(b []byte) {
	pl := p.extendPacketBody(len(b))
	copy(pl, b)
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

func (p *Packet) GetPacketID() uint32 {
	// _ = p.bytes.Bytes()[9]
	return *(*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
}

func (p *Packet) SetPacketID(id uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
	*pplen = id
}

// PacketCap  returns the current packetBody capacity
func (p *Packet) PacketCap() uint32 {
	return uint32(len(p.bytes.Bytes()))
}

// Release releases the packet to packet pool
func (p *Packet) Release() {

	bytespool.Put(p.bytes) // reclaim the buffer

	packetPool.Put(p)
}

func (p *Packet) reset() {
	p.src = nil
	p.setPacketBodyLen(0)
	p.SetPacketType(0)
	p.SetPacketFlags(0)
	p.SetPacketID(0)
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
