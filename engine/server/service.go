package server

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
// Write methods for a given Message will be called serially.

type Service interface {
	ServiceName

	ListenAndServe(network NetWorkName, address string, ops ServerOptions) error

	Close(chan struct{}) error
}