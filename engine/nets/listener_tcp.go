package nets

import (
	"errors"
	"net"
)

type TCPListener struct {
	network NetWorkName
}

func NewTCPListener(net NetWorkName) *TCPListener {
	return &TCPListener{net}
}

func (tl *TCPListener) Name() NetWorkName {
	return tl.network
}

func (tl *TCPListener) Listen(network NetWorkName, address string, ops *Options) (net.Listener, error) {
	if network != TCP && network != TCP6 {
		return nil, errors.New("not tcp")
	}

	return net.Listen(string(network), address)
}
