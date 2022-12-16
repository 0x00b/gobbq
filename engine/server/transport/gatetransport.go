package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type GateTransport struct {
	*ServerTransport
}

func NewGateTransport(ctx context.Context, conn net.Conn) *GateTransport {
	st := &GateTransport{}
	st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *GateTransport) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	// send to dispather
	return nil
}
