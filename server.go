package gobbq

import (
	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/engine/server/stream"
	"github.com/0x00b/gobbq/engine/server/stream/kcp"
	"github.com/0x00b/gobbq/engine/server/stream/tcp"
	"github.com/0x00b/gobbq/engine/server/stream/websocket"
)

// NewSever return gobbq server
func NewSever(opts ...server.ServerOption) *server.Server {
	svr := server.NewServer()

	svr.RegisterTransport(&websocket.WebSocket{})
	svr.RegisterTransport(stream.NewStreamTransport(tcp.NewTCPListener(server.TCP)))
	svr.RegisterTransport(stream.NewStreamTransport(tcp.NewTCPListener(server.TCP6)))
	svr.RegisterTransport(stream.NewStreamTransport(&kcp.KCPListener{}))

	return svr
}
