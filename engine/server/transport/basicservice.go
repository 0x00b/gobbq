package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
)

// type Transport interface {
// 	Serve() error
// }

type service struct {
	ctx         context.Context
	idleTimeout time.Duration
	lastVisited time.Time

	ops *server.ServerOptions

	// st Transport
}

func NewService(ctx context.Context) *service {
	return &service{
		ctx: ctx,
		// network:  lis.Name(),
	}
}

func (t *service) ListenAndServe(network server.NetWorkName, address string, opts *server.ServerOptions) error {
	t.ops = opts
	return t.listenAndServe(network, address, opts)
}

func (t *service) Close(chan struct{}) error {
	return nil
}

func (t *service) Name() string {
	return ""
}

// ===== inner =====

func (t *service) listenAndServe(network server.NetWorkName, address string, opts *server.ServerOptions) error {

	if network == server.WebSocket {
		return newWebSocketService().ListenAndServe(network, address, opts)
	}

	var ln net.Listener
	var err error

	switch network {
	case server.KCP:
		ln, err = NewDefaultKCPListener().Listen(network, address, opts)
	case server.TCP, server.TCP6:
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

func (t *service) handleConn(rawConn net.Conn, opts *server.ServerOptions) {
	if t.ops.TLSCertFile != "" && t.ops.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(t.ops.TLSCertFile, t.ops.TLSKeyFile)
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

	// t.st.Serve()
	// NewServerTransport(context.TODO(), conn).Serve()

	conn := &conn{
		rwc:              rawConn,
		ctx:              context.Background(),
		packetReadWriter: codec.NewPacketReadWriter(context.Background(), rawConn),
		PacketHandler:    NewServerPacketHandler(context.Background(), rawConn, opts),
		opts:             opts,
	}
	conn.Serve()

	return
}
