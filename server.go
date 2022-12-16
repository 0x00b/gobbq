package gobbq

import (
	"context"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/engine/server/transport"
	"github.com/0x00b/gobbq/engine/server/transport/kcp"
	"github.com/0x00b/gobbq/engine/server/transport/tcp"
	"github.com/0x00b/gobbq/engine/server/transport/websocket"
)

// NewSever return gobbq server
func NewSever(opts ...server.ServerOption) *server.Server {
	svr := server.NewServer()
	svr.RegisterTransport(&websocket.WebSocket{})
	svr.RegisterTransport(transport.NewTransport(context.Background(), tcp.NewTCPListener(server.TCP)))
	svr.RegisterTransport(transport.NewTransport(context.Background(), tcp.NewTCPListener(server.TCP6)))
	svr.RegisterTransport(transport.NewTransport(context.Background(), &kcp.KCPListener{}))

	return svr
}
