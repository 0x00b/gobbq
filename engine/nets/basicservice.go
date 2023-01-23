package nets

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

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

func (t *service) ListenAndServe(network NetWorkName, address string, opts *Options) error {
	t.opts = opts
	return t.listenAndServe(network, address, opts)
}

func (t *service) Close(chan struct{}) error {
	return nil
}

func (t *service) Name() string {
	return "service"
}

// ===== inner =====

func (t *service) listenAndServe(network NetWorkName, address string, opts *Options) error {

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
		go t.handleConn(conn, opts)
	}
}

func (t *service) handleConn(rawConn net.Conn, opts *Options) {
	if t.opts.TLSCertFile != "" && t.opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(t.opts.TLSCertFile, t.opts.TLSKeyFile)
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
		ctx:              context.Background(),
		packetReadWriter: codec.NewPacketReadWriter(context.Background(), rawConn),
		PacketHandler:    opts.PacketHandler,
		opts:             opts,
	}
	conn.Serve()

	return
}
