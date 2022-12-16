package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/server"
)

type ServerTransportInterface interface {
	Serve() error
}

type transport struct {
	ctx         context.Context
	idleTimeout time.Duration
	lastVisited time.Time

	network  server.NetWorkName
	listener Listener
	ops      server.ServerOptions

	st ServerTransportInterface
}

func NewTransport(ctx context.Context, lis Listener) *transport {
	return &transport{
		ctx:      ctx,
		listener: lis,
		network:  lis.Name(),
	}
}

type Listener interface {
	server.NetNamer
	Listen(network server.NetWorkName, address string, ops server.ServerOptions) (net.Listener, error)
}

func (t *transport) ListenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {
	t.ops = ops
	return t.listenAndServe(network, address, ops)
}
func (t *transport) Close(chan struct{}) error {
	return nil
}

func (t *transport) Name() server.NetWorkName {
	return t.network
}

// ===== inner =====

func (t *transport) listenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {

	ln, err := t.listener.Listen(network, address, ops)
	fmt.Printf("Listening on %s: %s ...", network, address)

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
		go t.handleConn(conn)
	}
}

func (t *transport) handleConn(conn net.Conn) {
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
		tlsConn := tls.Server(conn, tlsConfig)
		conn = net.Conn(tlsConn)
	}

	fmt.Println("handleconn")

	// t.st.Serve()
	NewServerTransport(context.TODO(), conn).Serve()

	return
}
