package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type DispatherPacketHandler struct {
}

func NewDispatherTransport(ctx context.Context, conn net.Conn) *DispatherPacketHandler {
	st := &DispatherPacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *DispatherPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))
	// send to game
	// or send to gate

	return nil
}
