package nets

import (
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type Options struct {
	Network     string
	Address     string
	CACertFile  string // ca证书
	TLSCertFile string // server证书
	TLSKeyFile  string // server秘钥

	maxSendPacketSize int
	writeBufferSize   int
	readBufferSize    int
	numServerWorkers  uint32
	connectionTimeout time.Duration

	PacketHandler PacketHandler
	ConnHandler   ConnHandler
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
type Option interface {
	apply(*Options)
}

type withPacketHandler struct {
	PacketHandler
}

func WithPacketHandler(ph PacketHandler) *withPacketHandler {
	return &withPacketHandler{ph}
}

func (w *withPacketHandler) apply(s *Options) {
	s.PacketHandler = w.PacketHandler
}

type withConnHandler struct {
	ConnHandler
}

func WithConnHandler(ph ConnHandler) *withConnHandler {
	return &withConnHandler{ph}
}

func (w *withConnHandler) apply(s *Options) {
	s.ConnHandler = w.ConnHandler
}
