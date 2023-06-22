package nets

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/0x00b/gobbq/engine/bytespool"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
)

// NOTE 如果需要通过chan 或者其他方式给其他协程使用，一定要retain，release
// NOTE 如果需要通过chan 或者其他方式给其他协程使用，一定要retain，release
// NOTE 如果需要通过chan 或者其他方式给其他协程使用，一定要retain，release
// Packet is a packet for sending data
type Packet struct {
	refcount int32

	Src *Conn // not nil indicates this is request packet

	// header: 只能在packet的生命周期内使用
	Header *bbq.Header

	totalLen  uint32
	headerLen uint32

	// data(header + body)
	bytes *bytespool.Bytes

	// 管理pkt的生命周期
	ctx context.Context
}

const (
	minPacketBufferLen  = bytespool.MinBufferCap
	MaxPacketBodyLength = bytespool.MaxBufferCap
)

var (
	packetEndian = binary.LittleEndian
	packetPool   = &sync.Pool{
		New: func() any {
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
type ReleasePkt func()

func allocPacket() *Packet {
	pkt := packetPool.Get().(*Packet)

	pkt.reset()

	// xlog.Printf("pkt pool get: %d", unsafe.Pointer(pkt))

	return pkt
}

// NewPacket allocates a new packet
func NewPacket() (*Packet, ReleasePkt) {
	pkt := allocPacket()
	return pkt, func() {
		// xlog.Printf("release callback %d", unsafe.Pointer(pkt))
		pkt.Release()
	}
}

func (p *Packet) Context() context.Context {
	return p.ctx
}

func (p *Packet) String() string {
	return fmt.Sprintf("hdrlen:%d, totallen:%d, hdr[%s] body[%s]", p.headerLen, p.totalLen, p.Header.String(), p.PacketBody())
}

// 想持有pkt，需要自行Retain/Release
func (p *Packet) Retain() {

	refcount := atomic.AddInt32(&p.refcount, 1)
	_ = refcount
	// xlog.Printf("retain pkt:%d, %d", unsafe.Pointer(p), refcount)

}

// Release releases the packet to packet pool
func (p *Packet) Release() {

	refcount := atomic.AddInt32(&p.refcount, -1)

	// xlog.Printf("release pkt:%d, %d", unsafe.Pointer(p), refcount)

	if refcount == 0 {
		// xlog.Printf("release pkt:%d, %d", unsafe.Pointer(p), unsafe.Pointer(p.bytes))

		bytespool.Put(p.bytes)
		packetPool.Put(p)
	} else if refcount < 0 {
		panic(fmt.Errorf("releasing packet with refcount=%d", p.refcount))
	}
}

// WriteBytes appends slice of bytes to the end of packetBody
func (p *Packet) WriteBody(b []byte) error {

	header, err := codec.DefaultCodec.Marshal(p.Header)
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

func (p *Packet) Serialize() []byte {
	pdata := p.Data()
	return append(packetEndian.AppendUint32(nil, uint32(len(pdata))), pdata...)
}

// PacketBody returns the total packetBody of packet
func (p *Packet) PacketBody() []byte {
	if p.bytes == nil {
		return nil
	}
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
	p.bytes = nil
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
		panic(errPacketBodyTooLarge)
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

	// xlog.Printf("bytes pool get: %d %d", unsafe.Pointer(p), unsafe.Pointer(bs))

	// copy(bs.Bytes(), p.Data())
	oldBytes := p.bytes
	p.bytes = bs

	bytespool.Put(oldBytes)

	return p.Data()
}
