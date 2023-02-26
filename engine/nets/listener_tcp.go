package nets

import (
	"errors"
	"net"
)

type TCPListener struct {
	network NetWorkName

	listener net.Listener

	opts *Options
}

func NewTCPListener(net NetWorkName, opts *Options) *TCPListener {
	return &TCPListener{
		network: net,
		opts:    opts,
	}
}

func (tl *TCPListener) Name() NetWorkName {
	return tl.network
}

func (tl *TCPListener) Listen(network NetWorkName, address string) (net.Listener, error) {

	if network != TCP && network != TCP6 {
		return nil, errors.New("not tcp")
	}
	var err error
	tl.listener, err = net.Listen(string(network), address)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Accept waits for and returns the next connection to the listener.
func (tl *TCPListener) Accept() (net.Conn, error) {
	rwc, err := tl.listener.Accept()
	if err != nil {
		return nil, err
	}

	tconn, ok := rwc.(*net.TCPConn)
	if !ok {
		return nil, errors.New("bug: conn is not tcpconn type")
	}

	// todo 研究一下，会导致收包延迟了50ms
	// if err := tconn.SetNoDelay(tl.opts.NetNoDelay); err != nil {
	// 	return nil, err
	// }

	if err := tconn.SetKeepAlive(tl.opts.NetKeepAlive); err != nil {
		return nil, err
	}

	return tconn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (tl *TCPListener) Close() error {
	return tl.listener.Close()
}

// Addr returns the listener's network address.
func (tl *TCPListener) Addr() net.Addr {
	return tl.Addr()
}
