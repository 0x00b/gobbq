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

	hdr := pkt.GetHeader()

	if hdr.GetMethod() == "register_proxy_entity" {
		fmt.Println("register", string(hdr.GetSrcEntity().ID))
		RegisterEntity(entity.EntityID(hdr.GetSrcEntity().ID), pkt.Src)
		return nil
	}

	fmt.Println("recv", hdr.String())
	// send to game
	// or send to gate
	Proxy(entity.EntityID(hdr.DstEntity.ID), pkt)

	return nil
}
