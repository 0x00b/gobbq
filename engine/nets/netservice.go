package nets

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
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

	ListenAndServe(network NetWorkName, address string, opts *Options) error

	Close(chan struct{}) error
}

type service struct {
	ctx         context.Context
	idleTimeout time.Duration
	lastVisited time.Time

	opts *Options
}

func NewNetService(ctx context.Context) *service {
	return &service{
		ctx: ctx,
	}
}

func (s *service) ListenAndServe(network NetWorkName, address string, opts *Options) error {
	s.opts = opts
	return s.listenAndServe(network, address, opts)
}

func (s *service) Close(chan struct{}) error {
	return nil
}

func (s *service) Name() string {
	return "service"
}

// ===== inner =====

func (s *service) listenAndServe(network NetWorkName, address string, opts *Options) error {

	if network == WebSocket {
		return newWebSocketService(s.ctx).ListenAndServe(network, address, opts)
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

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err != nil {
				continue
			} else {
				return err
			}
		}

		fmt.Printf("Connection from: %s", conn.RemoteAddr())
		go s.handleConn(conn, opts)
	}
}

func (s *service) handleConn(rawConn net.Conn, opts *Options) {
	if s.opts.TLSCertFile != "" && s.opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.opts.TLSCertFile, s.opts.TLSKeyFile)
		if err != nil {
			fmt.Println(err, "load RSA key & certificate failed")
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

	fmt.Println("handleconn")

	conn := &conn{
		rwc:              rawConn,
		ctx:              s.ctx,
		packetReadWriter: codec.NewPacketReadWriter(s.ctx, rawConn),
		PacketHandler:    opts.PacketHandler,
		opts:             opts,
	}
	conn.Serve()

	return
}