package server

import (
	"context"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type ServerOptions struct {
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
	HandlePacket(c context.Context, opts *ServerOptions, pkt *codec.Packet) error
}

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type ServerOption interface {
	apply(*ServerOptions)
}

type withPacketHandler struct {
	PacketHandler
}

func WithPacketHandler(ph PacketHandler) *withPacketHandler {
	return &withPacketHandler{ph}
}

func (w *withPacketHandler) apply(s *ServerOptions) {
	s.PacketHandler = w.PacketHandler
}
