package main

import (
	"context"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ProxyPacketHandler{}

type ProxyPacketHandler struct {
	*entity.MethodPacketHandler
}

func NewProxyPacketHandler() *ProxyPacketHandler {
	st := &ProxyPacketHandler{entity.NewMethodPacketHandler()}
	return st
}

func (st *ProxyPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	err := st.MethodPacketHandler.HandlePacket(c, pkt)
	if err == nil {
		// handle succ
		return nil
	}

	if entity.NotMyMethod(err) {
		hdr := pkt.Header
		// send to game
		// or send to gate
		Proxy(entity.EntityID(hdr.DstEntity.ID), pkt)
	}

	return err
}
