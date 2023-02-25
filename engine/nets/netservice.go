package nets

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

type NetWorkName string

const (
	WebSocket NetWorkName = "websocket"
	TCP       NetWorkName = "tcp"
	TCP6      NetWorkName = "tcp6"
	KCP       NetWorkName = "kcp"
)

type NetName interface {
	// Name returns the name of the Transport implementation.
	// the result cannot change between calls.
	Name() NetWorkName
}

type ServiceName interface {
	// Name returns the name of the Transport implementation.
	// the result cannot change between calls.
	Name() string
}

// ServerTransport is the common interface for all gRPC server-side transport
// implementations.
//
// Methods may be called concurrently from multiple goroutines, but
// Write methods for a given Packet will be called serially.

type NetService interface {
	ServiceName

	ListenAndServe() error

	Close(chan struct{})
}

type service struct {
	opts *Options

	ctx    context.Context
	cancel func()

	idleTimeout time.Duration
	lastVisited time.Time
}

func NewNetService(opts ...Option) *service {
	svc := &service{
		opts: &Options{},
	}

	svc.ctx, svc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(svc.opts)
	}

	return svc
}

func (s *service) ListenAndServe() error {
	return s.listenAndServe(s.opts.network, s.opts.address, s.opts)
}

func (s *service) Close(closeChan chan struct{}) {

	xlog.Infoln("server closing", s.opts.network, s.opts.address)

	s.cancel()

	closeChan <- struct{}{}

	return
}

func (s *service) Name() string {
	return "service"
}

// ===== inner =====

func (s *service) listenAndServe(network NetWorkName, address string, opts *Options) error {

	if network == WebSocket {
		return newWebSocketService().ListenAndServe(network, address, opts)
	}

	var ln net.Listener
	var err error

	switch network {
	case KCP:
		ln, err = NewDefaultKCPListener().Listen(network, address, opts)
	case TCP, TCP6:
		ln, err = NewTCPListener(network).Listen(network, address, opts)
	default:
		panic(fmt.Sprintf("unkown network:%s", network))
	}

	if err != nil {
		return err
	}
	xlog.Infoln("listenAndServe from:", network, address)

	defer ln.Close()

	for {
		select {
		case <-s.ctx.Done():
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			if err != nil {
				continue
			} else {
				return err
			}
		}

		xlog.Infof("Connection from: %s", conn.RemoteAddr())

		go s.handleConn(conn, opts)
	}
}

func (s *service) handleConn(rawConn net.Conn, opts *Options) {

	// tcpConn.SetNoDelay(consts.CLIENT_PROXY_SET_TCP_NO_DELAY)

	// kcpCon.SetReadBuffer(consts.CLIENT_PROXY_READ_BUFFER_SIZE)
	// conn.SetWriteBuffer(consts.CLIENT_PROXY_WRITE_BUFFER_SIZE)
	// // turn on turbo mode according to https://github.com/skywind3000/kcp/blob/master/README.en.md#protocol-configuration
	// conn.SetNoDelay(consts.KCP_NO_DELAY, consts.KCP_INTERNAL_UPDATE_TIMER_INTERVAL, consts.KCP_ENABLE_FAST_RESEND, consts.KCP_DISABLE_CONGESTION_CONTROL)
	// conn.SetStreamMode(consts.KCP_SET_STREAM_MODE)
	// conn.SetWriteDelay(consts.KCP_SET_WRITE_DELAY)
	// conn.SetACKNoDelay(consts.KCP_SET_ACK_NO_DELAY)

	if s.opts.TLSCertFile != "" && s.opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.opts.TLSCertFile, s.opts.TLSKeyFile)
		if err != nil {
			xlog.Traceln(err, "load RSA key & certificate failed")
			return
		}
		tlsConfig := &tls.Config{
			//MinVersion:       tls.VersionTLS12,
			//CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			Certificates: []tls.Certificate{cert},
			//CipherSuites: []uint16{
			//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			//	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			//},
			//PreferServerCipherSuites: true,
		}
		tlsConn := tls.Server(rawConn, tlsConfig)
		rawConn = net.Conn(tlsConn)
	}

	xlog.Traceln("handleconn")

	ctx, cancel := context.WithCancel(context.Background())
	cn := newDefaultConn(ctx)
	cn.cancel = cancel

	cn.rwc = rawConn
	cn.packetReadWriter = codec.NewPacketReadWriter(rawConn)
	cn.PacketHandler = opts.PacketHandler
	cn.opts = opts
	if opts.ConnHandler != nil {
		cn.ConnHandler = opts.ConnHandler
	}

	cn.Serve()

	return
}
