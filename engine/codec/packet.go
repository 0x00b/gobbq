package codec

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"
)

var (
	packetEndian                  = binary.LittleEndian
	predefinePacketBodyCapacities []uint32

	packetBufferPools = map[uint32]*sync.Pool{}
	packetPool        = &sync.Pool{
		New: func() interface{} {
			p := &Packet{}
			p.bytes = p.initialBytes[:]
			return p
		},
	}
)

const (
	minPacketBodyCap       = 128
	packetBodyCapGrowShift = uint(2)
	// (4 byte total len) + (1 byte msg type) + (1 byte flag) + (4 byte msg id)
	packetHeaderSize = 10
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

type Flags uint8

func init() {
	packetParsers[PacketRPC] = rpcPacketParser
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

// Packet-specific PacketHeader flag bits.
const (
	// Data Packet
	FlagDataZipCompress Flags = 0x1
)

func init() {
	packetBodyCap := uint32(minPacketBodyCap) << packetBodyCapGrowShift
	for packetBodyCap < MaxPacketBodyLength {
		predefinePacketBodyCapacities = append(predefinePacketBodyCapacities, packetBodyCap)
		packetBodyCap <<= packetBodyCapGrowShift
	}
	predefinePacketBodyCapacities = append(predefinePacketBodyCapacities, MaxPacketBodyLength)

	for _, packetBodyCap := range predefinePacketBodyCapacities {
		packetBodyCap := packetBodyCap
		packetBufferPools[packetBodyCap] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, packetHeaderSize+packetBodyCap)
			},
		}
	}
}

func getPacketBodyCapOfPacketBodyLen(packetBodyLen uint32) uint32 {
	for _, packetBodyCap := range predefinePacketBodyCapacities {
		if packetBodyCap >= packetBodyLen {
			return packetBodyCap
		}
	}
	return MaxPacketBodyLength
}

// Packet is a packet for sending data
type Packet struct {
	Src *PacketReadWriter

	// method       string // rpc mentod
	bytes        []byte
	initialBytes [packetHeaderSize + minPacketBodyCap]byte
}

func allocPacket() *Packet {
	pkt := packetPool.Get().(*Packet)

	if pkt.GetPacketBodyLen() != 0 {
		panic(fmt.Errorf("allocPacket: packetBody should be 0, but is %d", pkt.GetPacketBodyLen()))
	}

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
	return p.bytes[packetHeaderSize : packetHeaderSize+p.GetPacketBodyLen()]
}

// GetPacketBodyLen returns the packetBody length
func (p *Packet) GetPacketBodyLen() uint32 {
	_ = p.bytes[3]
	return *(*uint32)(unsafe.Pointer(&p.bytes[0]))
}

func (p *Packet) SetPacketBodyLen(plen uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes[0]))
	*pplen = plen
}

func (p *Packet) GetPacketType() PacketType {
	_ = p.bytes[4]
	return *(*PacketType)(unsafe.Pointer(&p.bytes[4]))
}

func (p *Packet) SetPacketType(typ PacketType) {
	pplen := (*PacketType)(unsafe.Pointer(&p.bytes[4]))
	*pplen = typ
}

func (p *Packet) GetPacketFlags() Flags {
	_ = p.bytes[5]
	return *(*Flags)(unsafe.Pointer(&p.bytes[5]))
}

func (p *Packet) SetPacketFlags(flags Flags) {
	pplen := (*Flags)(unsafe.Pointer(&p.bytes[5]))
	*pplen = flags
}

func (p *Packet) GetPacketID() uint32 {
	_ = p.bytes[9]
	return *(*uint32)(unsafe.Pointer(&p.bytes[6]))
}

func (p *Packet) SetPacketID(id uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes[6]))
	*pplen = id
}

// PacketBodyCap returns the current packetBody capacity
func (p *Packet) PacketBodyCap() uint32 {
	return uint32(len(p.bytes) - packetHeaderSize)
}

// Release releases the packet to packet pool
func (p *Packet) Release() {

	p.Src = nil

	packetBodyCap := p.PacketBodyCap()
	if packetBodyCap > minPacketBodyCap {
		buffer := p.bytes
		p.bytes = p.initialBytes[:]
		packetBufferPools[packetBodyCap].Put(buffer) // reclaim the buffer
	}

	p.SetPacketBodyLen(0)
	p.SetPacketType(0)
	p.SetPacketFlags(0)
	p.SetPacketID(0)

	packetPool.Put(p)
}

func (p *Packet) data() []byte {
	return p.bytes[0 : packetHeaderSize+p.GetPacketBodyLen()]
}
func (p *Packet) packetBodySlice(i, j uint32) []byte {
	return p.bytes[i+packetHeaderSize : j+packetHeaderSize]
}

// 返回的结果是在header的buf之后
func (p *Packet) extendPacketBody(size int) []byte {
	if size > MaxPacketBodyLength {
		panic(ErrPacketBodyTooLarge)
	}

	packetBodyLen := p.GetPacketBodyLen()
	newPacketBodyLen := packetBodyLen + uint32(size)
	oldCap := p.PacketBodyCap()

	if newPacketBodyLen <= oldCap { // most case
		p.SetPacketBodyLen(newPacketBodyLen)
		return p.packetBodySlice(packetBodyLen, newPacketBodyLen)
	}

	if newPacketBodyLen > MaxPacketBodyLength {
		panic(ErrPacketBodyTooLarge)
	}

	// try to find the proper capacity for the size bytes
	resizeToCap := getPacketBodyCapOfPacketBodyLen(newPacketBodyLen)

	buffer := packetBufferPools[resizeToCap].Get().([]byte)
	if len(buffer) != int(resizeToCap+packetHeaderSize) {
		panic(fmt.Errorf("buffer size should be %d, but is %d", resizeToCap, len(buffer)))
	}
	copy(buffer, p.data())
	oldBytes := p.bytes
	p.bytes = buffer

	if oldCap > minPacketBodyCap {
		// release old bytes
		packetBufferPools[oldCap].Put(oldBytes)
	}

	p.SetPacketBodyLen(newPacketBodyLen)
	return p.packetBodySlice(packetBodyLen, newPacketBodyLen)
}
