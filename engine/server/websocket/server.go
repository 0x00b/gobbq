package websocket

import (
	"net"

	"github.com/0x00b/gobbq/engine/server"
)

type WebSocketServerTransport struct {
	conn net.Conn
}

// Receive receives incoming packets using the given handler.
func (ws *WebSocketServerTransport) Receive(func(*server.Packet)) error {
	return nil
}

// Write sends the data for the given packet.
// Write may not be called on all packets.
func (ws *WebSocketServerTransport) Write(pkt *server.Packet) error {
	return nil
}

// WriteStatus sends the status of a packet to the client.  WriteStatus is
// the final call made on a packet and always occurs.
// WriteStatus(s *Packet, st *status.Status) error

// Close tears down the transport. Once it is called, the transport
// should not be accessed any more. All the pending packets and their
// handlers will be terminated asynchronously.
func (ws *WebSocketServerTransport) Close() {
	return
}

// Drain notifies the client this ServerTransport stops accepting new RPCs.
func (ws *WebSocketServerTransport) Drain() {
	return
}

// RemoteAddr returns the remote network address.
func (ws *WebSocketServerTransport) RemoteAddr() net.Addr {
	return nil
}

func (ws *WebSocketServerTransport) Name() server.NetWorkName {
	return server.WebSocket
}
