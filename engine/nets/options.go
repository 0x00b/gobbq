package nets

import (
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
)

type Options struct {
	network     NetWorkName
	address     string
	CACertFile  string // ca证书
	TLSCertFile string // server证书
	TLSKeyFile  string // server秘钥

	CompressType bbq.CompressType //压缩类型
	ContentType  bbq.ContentType  //协议编码类型
	CheckFlags   uint32

	maxSendPacketSize int
	writeBufferSize   int
	readBufferSize    int
	numServerWorkers  uint32
	connectionTimeout time.Duration
	requestTimeout    time.Duration

	PacketHandler PacketHandler
	ConnHandler   ConnHandler
}

var DefaultOptions = &Options{
	CACertFile:        "",
	TLSCertFile:       "",
	TLSKeyFile:        "",
	CompressType:      bbq.CompressType_None,
	ContentType:       bbq.ContentType_Proto,
	maxSendPacketSize: 0,
	writeBufferSize:   0,
	readBufferSize:    0,
	numServerWorkers:  0,
	connectionTimeout: 0,
	requestTimeout:    0,
	PacketHandler:     nil,
	ConnHandler:       &defaultConnHandler{},
}

type PacketHandler interface {
	HandlePacket(pkt *codec.Packet) error
}

type ConnHandler interface {
	HandleEOF(*codec.PacketReadWriter)
	HandleTimeOut(*codec.PacketReadWriter)
	HandleFail(*codec.PacketReadWriter)
}

type defaultConnHandler struct {
}

func (ch *defaultConnHandler) HandleEOF(prw *codec.PacketReadWriter)     {}
func (ch *defaultConnHandler) HandleTimeOut(prw *codec.PacketReadWriter) {}
func (ch *defaultConnHandler) HandleFail(prw *codec.PacketReadWriter)    {}

// A Option sets options such as credentials, codec and keepalive parameters, etc.
type Option func(*Options)

func WithPacketHandler(ph PacketHandler) Option {
	return func(o *Options) {
		o.PacketHandler = ph
	}
}

func WithConnHandler(ph ConnHandler) Option {
	return func(o *Options) {
		o.ConnHandler = ph
	}
}

func WithNetwork(nw NetWorkName) Option {
	return func(o *Options) {
		o.network = nw
	}
}
func WithAddress(ad string) Option {
	return func(o *Options) {
		o.address = ad
	}
}
func WithCACertFile(CACertFile string) Option {
	return func(o *Options) {
		o.CACertFile = CACertFile
	}
}
func WithTLSCertFile(TLSCertFile string) Option {
	return func(o *Options) {
		o.TLSCertFile = TLSCertFile
	}
}
func WithTLSKeyFile(TLSKeyFile string) Option {
	return func(o *Options) {
		o.TLSKeyFile = TLSKeyFile
	}
}

func WithCompressType(CompressType bbq.CompressType) Option {
	return func(o *Options) {
		o.CompressType = CompressType
	}
}
func WithContentType(ContentType bbq.ContentType) Option {
	return func(o *Options) {
		o.ContentType = ContentType
	}
}

func WithMaxSendPacketSize(maxSendPacketSize int) Option {
	return func(o *Options) {
		o.maxSendPacketSize = maxSendPacketSize
	}
}
func WithWriteBufferSize(writeBufferSize int) Option {
	return func(o *Options) {
		o.writeBufferSize = writeBufferSize
	}
}
func WithReadBufferSize(readBufferSize int) Option {
	return func(o *Options) {
		o.readBufferSize = readBufferSize
	}
}
func WithNumServerWorkers(numServerWorkers uint32) Option {
	return func(o *Options) {
		o.numServerWorkers = numServerWorkers
	}
}
func WithConnectionTimeout(connectionTimeout time.Duration) Option {
	return func(o *Options) {
		o.connectionTimeout = connectionTimeout
	}
}
func WithRequestTimeout(requestTimeout time.Duration) Option {
	return func(o *Options) {
		o.requestTimeout = requestTimeout
	}
}

func WithCheckFlags(CheckFlags uint32) Option {
	return func(o *Options) {
		o.CheckFlags |= CheckFlags
	}
}
