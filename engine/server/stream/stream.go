package stream

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/pkg/errors"
)

type StreamListener interface {
	server.Namer
	Listen(network server.NetWorkName, address string, ops server.ServerOptions) (net.Listener, error)
}

type StreamTransport struct {
	network  server.NetWorkName
	listener StreamListener
	ops      server.ServerOptions
}

func NewStreamTransport(lis StreamListener) *StreamTransport {
	st := &StreamTransport{
		listener: lis,
		network:  lis.Name(),
	}

	return st
}

func (ts *StreamTransport) ListenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {
	ts.ops = ops
	return ts.listenAndServe(network, address, ops)
}

func (ts *StreamTransport) Close(chan struct{}) error {
	return nil
}

func (ts *StreamTransport) Name() server.NetWorkName {
	return ts.network
}

// ===== inner =====

func (ts *StreamTransport) listenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {

	ln, err := ts.listener.Listen(network, address, ops)
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
		go ts.handleConn(conn)
	}
}

func (ts *StreamTransport) handleConn(conn net.Conn) {
	if ts.ops.TLSCertFile != "" && ts.ops.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(ts.ops.TLSCertFile, ts.ops.TLSKeyFile)
		if err != nil {
			fmt.Println(errors.Wrap(err, "load RSA key & certificate failed"))
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

	NewStreamServer(context.TODO(), conn).Serve()

	return
}

// type ServerTransport interface {
// 	server.Namer

// 	// Receive receives incoming packets using the given handler.
// 	Receive(func(*Packet)) error

// 	// Write sends the data for the given packet.
// 	// Write may not be called on all packets.
// 	Write(pkt *Packet) error

// 	// WriteStatus sends the status of a packet to the client.  WriteStatus is
// 	// the final call made on a packet and always occurs.
// 	// WriteStatus(s *Packet, st *status.Status) error

// 	// Close tears down the transport. Once it is called, the transport
// 	// should not be accessed any more. All the pending packets and their
// 	// handlers will be terminated asynchronously.
// 	Close()

// 	// Drain notifies the client this ServerTransport stops accepting new RPCs.
// 	Drain()

// 	// RemoteAddr returns the remote network address.
// 	RemoteAddr() net.Addr

// 	// WriteHeader sends the header metadata for the given packet.
// 	// WriteHeader may not be called on all packets.
// 	// WriteHeader(s *Packet, md metadata.MD) error

// }

// // ClientTransport is the common interface for all gRPC client-side transport
// // implementations.
// type ClientTransport interface {
// 	server.Namer

// 	// Close tears down this transport. Once it returns, the transport
// 	// should not be accessed any more. The caller must make sure this
// 	// is called only once.
// 	Close(err error)

// 	// Write sends the data for the given packet. A nil packet indicates
// 	// the write is to be performed on the transport as a whole.
// 	Write(pkt *Packet) error

// 	// NewPacket creates a Packet for a Server Call.
// 	NewPacket(ctx context.Context, callHdr *CallHdr) (*Packet, error)

// 	// Error returns a channel that is closed when some I/O error
// 	// happens. Typically the caller should have a goroutine to monitor
// 	// this in order to take action (e.g., close the current transport
// 	// and create a new one) in error case. It should not return nil
// 	// once the transport is initiated.
// 	Error() <-chan struct{}

// 	// RemoteAddr returns the remote network address.
// 	RemoteAddr() net.Addr
// }

// // CallHdr carries the information of a particular RPC.
// type CallHdr struct {
// 	// Host specifies the peer's host.
// 	Host string

// 	// Method specifies the operation to perform.
// 	Method string

// 	// SendCompress specifies the compression algorithm applied on
// 	// outbound message.
// 	SendCompress string

// 	// Creds specifies credentials.PerRPCCredentials for a call.
// 	// Creds credentials.PerRPCCredentials

// 	// ContentSubtype specifies the content-subtype for a request. For example, a
// 	// content-subtype of "proto" will result in a content-type of
// 	// "application/gobbq+proto". The value of ContentSubtype must be all
// 	// lowercase, otherwise the behavior is undefined.
// 	ContentSubtype string

// 	DoneFunc func() // called when the packet is finished
// }
