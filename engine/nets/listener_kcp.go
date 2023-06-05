package nets

import (
	"errors"
	"net"

	"github.com/xtaci/kcp-go"
)

type KCPListener struct {
	opts *Options

	listener *kcp.Listener
}

func NewDefaultKCPListener(opts *Options) *KCPListener {
	return &KCPListener{
		opts: opts,
	}
}

func (kl *KCPListener) Listen(network NetWorkName, address string) (net.Listener, error) {
	if network != KCP {
		return nil, errors.New("not kcp")
	}
	var err error

	kl.listener, err = kcp.ListenWithOptions(address, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	return kl, nil
}

func (kl *KCPListener) Name() NetWorkName {
	return KCP
}

// Accept waits for and returns the next connection to the listener.
func (kl *KCPListener) Accept() (net.Conn, error) {
	rwc, err := kl.listener.Accept()
	if err != nil {
		return nil, err
	}

	kconn, ok := rwc.(*kcp.UDPSession)
	if !ok {
		return nil, errors.New("bug: conn is not kcpconn type")
	}
	// kconn.SetReadBuffer(consts.CLIENT_PROXY_READ_BUFFER_SIZE)
	// kconn.SetWriteBuffer(consts.CLIENT_PROXY_WRITE_BUFFER_SIZE)

	// turn on turbo mode according to https://github.com/skywind3000/kcp/blob/master/README.en.md#protocol-configuration
	if kl.opts.NetNoDelay {
		kconn.SetNoDelay(1, kl.opts.KcpInternalUpdateTimerInterval, kl.opts.KcpEnableFastResend, kl.opts.KcpDisableCongestionControl)
	} else {
		kconn.SetNoDelay(0, kl.opts.KcpInternalUpdateTimerInterval, kl.opts.KcpEnableFastResend, kl.opts.KcpDisableCongestionControl)
	}

	kconn.SetStreamMode(kl.opts.KcpSetStreamMode)
	kconn.SetWriteDelay(kl.opts.KcpSetWriteDelay)
	kconn.SetACKNoDelay(kl.opts.KcpSetAckNoDelay)
	kconn.SetMtu(kl.opts.KcpMTU)

	return kconn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (kl *KCPListener) Close() error {
	return kl.listener.Close()
}

// Addr returns the listener's network address.
func (kl *KCPListener) Addr() net.Addr {
	return kl.listener.Addr()
}
