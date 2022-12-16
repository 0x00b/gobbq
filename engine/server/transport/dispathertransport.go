package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type DispatherTransport struct {
	*ServerTransport
}

func NewDispatherTransport(ctx context.Context, conn net.Conn) *DispatherTransport {
	st := &DispatherTransport{}
	st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *DispatherTransport) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))
	// send to game
	// or send to gate

	return nil
}
