package tcp

import (
	"errors"
	"net"

	"github.com/0x00b/gobbq/engine/server"
)

type TCPListener struct {
	network server.NetWorkName
}

func NewTCPListener(net server.NetWorkName) *TCPListener {
	return &TCPListener{net}
}

func (tl *TCPListener) Name() server.NetWorkName {
	return tl.network
}

func (tl *TCPListener) Listen(network server.NetWorkName, address string, ops server.ServerOptions) (net.Listener, error) {
	if network != server.TCP && network != server.TCP6 {
		return nil, errors.New("not tcp")
	}

	return net.Listen(string(network), address)
}
