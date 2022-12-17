package gobbq

import (
	"context"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/engine/server/transport"
)

// NewSever return gobbq server
func NewSever(opts ...server.ServerOption) *server.Server {
	svr := server.NewServer()
	svr.RegisterService(transport.NewService(context.Background()))

	return svr
}
