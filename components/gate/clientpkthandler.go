package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
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

	// new client
	if hdr.GetMethod() == "new_client" {
		err := ex.RegisterEntity(hdr.GetSrcEntity().ID)

		if err != nil {
			panic(err)
		}

		fmt.Println("register", string(hdr.GetSrcEntity().ID))

		RegisterEntity(entity.EntityID(hdr.GetSrcEntity().ID), pkt.Src)

		return err
	}

	// send to proxy
	return ex.SendProxy(pkt)

}
