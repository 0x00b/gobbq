package kcp

import (
	"errors"
	"net"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/xtaci/kcp-go"
)

type KCPListener struct {
}

func NewDefaultKCPListener() *KCPListener {
	return &KCPListener{}
}

func (kl *KCPListener) Listen(network server.NetWorkName, address string, ops server.ServerOptions) (net.Listener, error) {
	if network != server.KCP {
		return nil, errors.New("not kcp")
	}

	return kcp.ListenWithOptions(address, nil, 10, 3)
}

func (kl *KCPListener) Name() server.NetWorkName {
	return server.KCP
}
