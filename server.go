package gobbq

import (
	"context"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/engine/server/transport"
)

// NewSever return gobbq server
func NewSever(opts ...server.ServerOption) *server.Server {
	svr := server.NewServer(opts...)
	svr.RegisterNetService(transport.NewNetService(context.Background()))

	return svr
}
