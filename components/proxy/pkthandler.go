package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
)

var _ server.PacketHandler = &ProxyPacketHandler{}

type ProxyPacketHandler struct {
}

func NewProxyPacketHandler() *ProxyPacketHandler {
	st := &ProxyPacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *ProxyPacketHandler) HandlePacket(c context.Context, opts *server.ServerOptions, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))
	// send to game
	// or send to gate

	return nil
}
