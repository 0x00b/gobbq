package codec

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/0x00b/gobbq/engine/bytespool"
)

var (
	messageEndian = binary.LittleEndian
	messagePool   = &sync.Pool{
		New: func() interface{} {
			p := &Message{}
			p.bytes = bytespool.Get(minMessageBufferLen)
			return p
		},
	}
)

const (
	// (4 byte total len) + (1 byte msg type) + (1 byte flag) + (4 byte message id)
	messageHeaderSize    = 10
	minMessageBufferLen  = bytespool.MinBufferCap
	MaxMessageBodyLength = bytespool.MaxBufferCap - messageHeaderSize
)

type MessageType uint8

const (
	MessageRPC MessageType = 0x0
	MessageSys MessageType = 0x1
)

var messageName = map[MessageType]string{
	MessageRPC: "RPC",
	MessageSys: "PING",
}

func (t MessageType) String() string {
	if s, ok := messageName[t]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN_MESSAGE_TYPE_%d", uint8(t))
}

type messageParser func(pc *MessageReadWriter, pkt *Message) (uint32, error)

var messageParsers map[MessageType]messageParser

// depend on MessageType
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
	//MessageRPC
	FlagDataRpcResponse  Flags = 0x01 // 0x00：request， 0x01：response
	FlagDataZipCompress  Flags = 0x02 // 留两位来代表压缩算法 {0x2,0x4,0x6}
	FlagDataChecksumIEEE Flags = 0x08
	FlagDataProtoBuf     Flags = 0x10
)

var flagName = map[MessageType]map[Flags]string{
	MessageRPC: {
		FlagDataRpcResponse:  "RpcResponse",
		FlagDataZipCompress:  "ZipCompress",
		FlagDataChecksumIEEE: "ChecksumIEEE",
		FlagDataProtoBuf:     "ProtoBuf",
	},
}

func init() {
	messageParsers = map[MessageType]messageParser{
		MessageRPC: rpcMessageParser,
	}
}

func rpcMessageParser(pc *MessageReadWriter, pkt *Message) (uint32, error) {
	// var buff16uint [2]byte
	// io.ReadFull(pc.rw, buff16uint[:])
	// mothodLen := messageEndian.Uint16(buff16uint[:])
	// buf := make([]byte, mothodLen)
	// _, err := io.ReadFull(pc.rw, buf)
	// if err != nil {
	// 	return 0, err
	// }
	// pkt.method = string(buf)
	// return 2 + uint32(mothodLen), nil
	return 0, nil
}

// Message is a message for sending data
type Message struct {
	refcount int32

	Src *MessageReadWriter // not nil indicates this is request message

	// depend on MessageType
	// meta interface{}

	bytes *bytespool.Bytes

	// members:
	// MessageID uint32 <-> RequestID <-> ResponseID, represent one call.
	// MessageType
	// MessageFlags
}

func allocMessage() *Message {
	pkt := messagePool.Get().(*Message)

	pkt.reset()

	return pkt
}

// NewMessage allocates a new message
func NewMessage() *Message {
	return allocMessage()
}

func (p *Message) Retain() {
	atomic.AddInt32(&p.refcount, 1)
}

// Release releases the message to message pool
func (p *Message) Release() {

	refcount := atomic.AddInt32(&p.refcount, -1)

	if refcount == 0 {
		bytespool.Put(p.bytes)
		messagePool.Put(p)
	} else if refcount < 0 {
		panic(fmt.Errorf("releasing message with refcount=%d", p.refcount))
	}
}

// WriteBytes appends slice of bytes to the end of messageBody
func (p *Message) WriteBytes(b []byte) {
	pl := p.extendMessageBody(len(b))
	copy(pl, b)
}

// MessageBody returns the total messageBody of message
func (p *Message) MessageBody() []byte {
	return p.bytes.Bytes()[messageHeaderSize : messageHeaderSize+p.GetMessageBodyLen()]
}

// GetMessageBodyLen returns the messageBody length
func (p *Message) GetMessageBodyLen() uint32 {
	// _ = p.bytes.Bytes()[3]
	return *(*uint32)(unsafe.Pointer(&p.bytes.Bytes()[0]))
}

func (p *Message) setMessageBodyLen(plen uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes.Bytes()[0]))
	*pplen = plen
}

func (p *Message) GetMessageType() MessageType {
	// _ = p.bytes.Bytes()[4]
	return *(*MessageType)(unsafe.Pointer(&p.bytes.Bytes()[4]))
}

func (p *Message) SetMessageType(typ MessageType) {
	pplen := (*MessageType)(unsafe.Pointer(&p.bytes.Bytes()[4]))
	*pplen = typ
}

func (p *Message) GetMessageFlags() Flags {
	// _ = p.bytes.Bytes()[5]
	return *(*Flags)(unsafe.Pointer(&p.bytes.Bytes()[5]))
}

func (p *Message) SetMessageFlags(flags Flags) {
	pplen := (*Flags)(unsafe.Pointer(&p.bytes.Bytes()[5]))
	*pplen = flags
}

// MessageID <-> RequestID <-> ResponseID, represent one call.
func (p *Message) GetMessageID() uint32 {
	// _ = p.bytes.Bytes()[9]
	return *(*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
}

// MessageID <-> RequestID <-> ResponseID, represent one call.
func (p *Message) SetMessageID(id uint32) {
	pplen := (*uint32)(unsafe.Pointer(&p.bytes.Bytes()[6]))
	*pplen = id
}

// MessageCap  returns the current messageBody capacity
func (p *Message) MessageCap() uint32 {
	return uint32(len(p.bytes.Bytes()))
}

func (p *Message) reset() {
	p.Src = nil
	p.setMessageBodyLen(0)
	p.SetMessageType(0)
	p.SetMessageFlags(0)
	p.SetMessageID(0)
	p.refcount = 1
}

func (p *Message) Data() []byte {
	return p.bytes.Bytes()[0 : messageHeaderSize+p.GetMessageBodyLen()]
}

func (p *Message) messageBodySlice(i, j uint32) []byte {
	return p.bytes.Bytes()[i+messageHeaderSize : j+messageHeaderSize]
}

// 返回的结果是在header的buf之后
func (p *Message) extendMessageBody(size int) []byte {
	if size > MaxMessageBodyLength {
		panic(ErrMessageBodyTooLarge)
	}

	messageBodyLen := p.GetMessageBodyLen()
	newMessageBodyLen := messageBodyLen + uint32(size)
	newMessageLen := messageHeaderSize + newMessageBodyLen
	oldCap := p.MessageCap()

	if newMessageLen <= oldCap { // most case
		p.setMessageBodyLen(newMessageBodyLen)
		return p.messageBodySlice(messageBodyLen, newMessageBodyLen)
	}

	// get new buffer

	if newMessageLen > MaxMessageBodyLength {
		panic(ErrMessageBodyTooLarge)
	}
	bs := bytespool.Get(newMessageLen)

	copy(bs.Bytes(), p.Data())
	oldBytes := p.bytes
	p.bytes = bs

	bytespool.Put(oldBytes)

	p.setMessageBodyLen(newMessageBodyLen)
	return p.messageBodySlice(messageBodyLen, newMessageBodyLen)
}
