package gobbq

import (
	"github.com/0x00b/gobbq/engine/nets"
)

// NewSever return gobbq server
func NewSever(opts ...nets.Option) *nets.Server {
	svr := nets.NewServer(opts...)
	svr.RegisterNetService(nets.NewNetService())

	return svr
}
