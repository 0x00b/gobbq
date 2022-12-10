package server

import (
	"context"
	"net"
)

type NetWorkName string

const (
	WebSocket NetWorkName = "websocket"
	TCP       NetWorkName = "tcp"
	TCP6      NetWorkName = "tcp6"
	KCP       NetWorkName = "kcp"
)

var registeredTransport = make(map[NetWorkName]Transport)

func RegisterTransport(t Transport) {
	registeredTransport[t.Name()] = t
}

type Namer interface {
	// Name returns the name of the Transport implementation.
	// the result cannot change between calls.
	Name() NetWorkName
}

// ServerTransport is the common interface for all gRPC server-side transport
// implementations.
//
// Methods may be called concurrently from multiple goroutines, but
// Write methods for a given Packet will be called serially.

type Transport interface {
	Namer

	ListenAndServe(network NetWorkName, address string, ops ServerOptions) error

	Close(chan struct{}) error
}

type ServerTransport interface {
	Namer

	// Receive receives incoming packets using the given handler.
	Receive(func(*Packet)) error

	// Write sends the data for the given packet.
	// Write may not be called on all packets.
	Write(pkt *Packet) error

	// WriteStatus sends the status of a packet to the client.  WriteStatus is
	// the final call made on a packet and always occurs.
	// WriteStatus(s *Packet, st *status.Status) error

	// Close tears down the transport. Once it is called, the transport
	// should not be accessed any more. All the pending packets and their
	// handlers will be terminated asynchronously.
	Close()

	// Drain notifies the client this ServerTransport stops accepting new RPCs.
	Drain()

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr

	// WriteHeader sends the header metadata for the given packet.
	// WriteHeader may not be called on all packets.
	// WriteHeader(s *Packet, md metadata.MD) error

}

// ClientTransport is the common interface for all gRPC client-side transport
// implementations.
type ClientTransport interface {
	Namer

	// Close tears down this transport. Once it returns, the transport
	// should not be accessed any more. The caller must make sure this
	// is called only once.
	Close(err error)

	// Write sends the data for the given packet. A nil packet indicates
	// the write is to be performed on the transport as a whole.
	Write(pkt *Packet) error

	// NewPacket creates a Packet for a Server Call.
	NewPacket(ctx context.Context, callHdr *CallHdr) (*Packet, error)

	// Error returns a channel that is closed when some I/O error
	// happens. Typically the caller should have a goroutine to monitor
	// this in order to take action (e.g., close the current transport
	// and create a new one) in error case. It should not return nil
	// once the transport is initiated.
	Error() <-chan struct{}

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}

// CallHdr carries the information of a particular RPC.
type CallHdr struct {
	// Host specifies the peer's host.
	Host string

	// Method specifies the operation to perform.
	Method string

	// SendCompress specifies the compression algorithm applied on
	// outbound message.
	SendCompress string

	// Creds specifies credentials.PerRPCCredentials for a call.
	// Creds credentials.PerRPCCredentials

	// ContentSubtype specifies the content-subtype for a request. For example, a
	// content-subtype of "proto" will result in a content-type of
	// "application/gobbq+proto". The value of ContentSubtype must be all
	// lowercase, otherwise the behavior is undefined.
	ContentSubtype string

	DoneFunc func() // called when the packet is finished
}
