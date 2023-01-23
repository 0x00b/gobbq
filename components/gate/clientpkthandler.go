package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ClientPacketHandler{}

type ClientPacketHandler struct {
}

func NewClientPacketHandler() *ClientPacketHandler {
	st := &ClientPacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *ClientPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	hdr := pkt.GetHeader()

	_ = hdr

	// send to proxy

	return nil
}
