package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type GatePacketHandler struct {
}

func NewGatePacketHandler(ctx context.Context, conn net.Conn) *GatePacketHandler {
	st := &GatePacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *GatePacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	// send to dispather
	return nil
}
