package nets

import (
	"context"
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
}

type PacketHandler interface {
	HandlePacket(c context.Context, pkt *codec.Packet) error
}

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
