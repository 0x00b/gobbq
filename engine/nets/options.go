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
	connectionTimeout time.Duration
	RequestTimeout    time.Duration
	MaxCloseWaitTime  time.Duration

	NetNoDelay bool

	// for tcp
	NetKeepAlive bool

	// for kcp
	KcpSetStreamMode bool
	KcpSetWriteDelay bool

	KcpSetAckNoDelay bool

	// for kcp nodelay
	KcpInternalUpdateTimerInterval int
	KcpEnableFastResend            int
	KcpDisableCongestionControl    int

	// for kcp mtu
	KcpMTU int

	PacketHandler  PacketHandler
	ConnErrHandler ConnErrHandler
}

var DefaultOptions = &Options{
	network:                        TCP,
	address:                        "",
	CACertFile:                     "",
	TLSCertFile:                    "",
	TLSKeyFile:                     "",
	CompressType:                   bbq.CompressType_None,
	ContentType:                    bbq.ContentType_Proto,
	CheckFlags:                     0,
	maxSendPacketSize:              0,
	writeBufferSize:                0,
	readBufferSize:                 0,
	connectionTimeout:              5 * time.Second,
	RequestTimeout:                 5 * time.Second,
	MaxCloseWaitTime:               10,
	NetNoDelay:                     true,
	NetKeepAlive:                   true,
	KcpSetStreamMode:               true,
	KcpSetWriteDelay:               true,
	KcpSetAckNoDelay:               true,
	KcpInternalUpdateTimerInterval: 10,
	KcpEnableFastResend:            2,
	KcpDisableCongestionControl:    1,
	KcpMTU:                         1400,
	PacketHandler:                  nil,
	ConnErrHandler:                 nil,
}

type PacketHandler interface {
	HandlePacket(pkt *codec.Packet) error
}

type ConnErrHandler interface {
	HandleEOF(*conn)
	HandleTimeOut(*conn)
	HandleFail(*conn)
}

// A Option sets options such as credentials, codec and keepalive parameters, etc.
type Option func(*Options)

func WithPacketHandler(ph PacketHandler) Option {
	return func(o *Options) {
		o.PacketHandler = ph
	}
}

func WithConnErrHandler(ph ConnErrHandler) Option {
	return func(o *Options) {
		o.ConnErrHandler = ph
	}
}

func WithNetwork(network NetWorkName, address string) Option {
	return func(o *Options) {
		o.network = network
		o.address = address
	}
}

func WithTls(CACertFile, TLSCertFile, TLSKeyFile string) Option {
	return func(o *Options) {
		o.CACertFile = CACertFile
		o.TLSCertFile = TLSCertFile
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

func WithConnectionTimeout(connectionTimeout time.Duration) Option {
	return func(o *Options) {
		o.connectionTimeout = connectionTimeout
	}
}
func WithRequestTimeout(requestTimeout time.Duration) Option {
	return func(o *Options) {
		o.RequestTimeout = requestTimeout
	}
}

func WithCheckFlags(CheckFlags uint32) Option {
	return func(o *Options) {
		o.CheckFlags |= CheckFlags
	}
}

func WithNetNoDelay(NetNoDelay bool) Option {
	return func(o *Options) {
		o.NetNoDelay = NetNoDelay
	}
}

func WithTcpKeepAlive(NetKeepAlive bool) Option {
	return func(o *Options) {
		o.NetKeepAlive = NetKeepAlive
	}
}

func WithKcpSetStreamMode(KcpSetStreamMode bool) Option {
	return func(o *Options) {
		o.KcpSetStreamMode = KcpSetStreamMode
	}
}

func WithKcpSetWriteDelay(KcpSetWriteDelay bool) Option {
	return func(o *Options) {
		o.KcpSetWriteDelay = KcpSetWriteDelay
	}
}

func WithKcpSetAckNoDelay(KcpSetAckNoDelay bool) Option {
	return func(o *Options) {
		o.KcpSetAckNoDelay = KcpSetAckNoDelay
	}
}

func WithKcpNodelay(nodelay bool, intval, fr, nc int) Option {
	return func(o *Options) {
		o.NetNoDelay = nodelay
		o.KcpInternalUpdateTimerInterval = intval
		o.KcpEnableFastResend = fr
		o.KcpDisableCongestionControl = nc
	}
}

func WithKcpMTU(mtu int) Option {
	return func(o *Options) {
		o.KcpMTU = mtu
	}
}
