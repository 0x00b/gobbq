package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ProxyPacketHandler{}

type ProxyPacketHandler struct {
}

func NewProxyPacketHandler() *ProxyPacketHandler {
	st := &ProxyPacketHandler{}
	return st
}

func (st *ProxyPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	hdr := pkt.GetHeader()

	id := hdr.GetDstEntity().ID
	rw, ok := cltMap[entity.EntityID(id)]
	if !ok {
		fmt.Println("unknown client")
		return nil
	}

	return rw.WritePacket(pkt)

}
