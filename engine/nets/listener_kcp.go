package nets

import (
	"errors"
	"net"

	"github.com/xtaci/kcp-go"
)

type KCPListener struct {
}

func NewDefaultKCPListener() *KCPListener {
	return &KCPListener{}
}

func (kl *KCPListener) Listen(network NetWorkName, address string, ops *Options) (net.Listener, error) {
	if network != KCP {
		return nil, errors.New("not kcp")
	}

	return kcp.ListenWithOptions(address, nil, 10, 3)
}

func (kl *KCPListener) Name() NetWorkName {
	return KCP
}
