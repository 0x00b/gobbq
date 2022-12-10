package server

import (
	"context"

	"github.com/0x00b/gobbq/engine/entity"
)

// Packet represents an RPC in the transport layer.
type Packet struct {
	id           uint32
	st           ServerTransport    // nil for client side Stream
	ct           ClientTransport    // nil for server side Stream
	ctx          context.Context    // the associated context of the stream
	cancel       context.CancelFunc // always nil for client side Stream
	done         chan struct{}      // closed at the end of stream to unblock writers. On the client side.
	doneFunc     func()             // invoked at the end of stream on client side.
	ctxDone      <-chan struct{}    // same as done chan but for server side. Cache of ctx.Done() (for performance)
	method       string             // the associated RPC method of the stream
	recvCompress string
	sendCompress string
	// buf          *recvBuffer
	// trReader     io.Reader

	bytesReceived uint32 // indicates whether any bytes have been received on this stream

	// contentSubtype is the content-subtype for requests.
	// this must be lowercase or the behavior is undefined.
	contentSubtype string

	entityID entity.EntityID
}

func (p *Packet) Method() string {
	return p.method
}

func (p *Packet) EntityID() entity.EntityID {
	return p.entityID
}
