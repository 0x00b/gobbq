package gobbq

import (
	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/engine/server/kcp"
	"github.com/0x00b/gobbq/engine/server/tcp"
	"github.com/0x00b/gobbq/engine/server/websocket"
)

func init() {
	server.RegisterTransport(&websocket.WebSocket{})
	server.RegisterTransport(tcp.NewTCPTransport(server.TCP))
	server.RegisterTransport(tcp.NewTCPTransport(server.TCP6))
	server.RegisterTransport(&kcp.KCPTransport{})
}

// NewSever return gobbq server
func NewSever(opts ...server.ServerOption) *server.Server {
	svr := &server.Server{}

	return svr
}
