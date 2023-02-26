package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ClientPacketHandler{}

type ClientPacketHandler struct {
	nets.PacketHandler
}

func NewClientPacketHandler(etyMgr *entity.EntityManager) *ClientPacketHandler {
	st := &ClientPacketHandler{
		PacketHandler: etyMgr,
	}
	return st
}

func (st *ClientPacketHandler) HandlePacket(pkt *codec.Packet) error {

	err := st.PacketHandler.HandlePacket(pkt)
	if err == nil {
		// handle succ
		return nil
	}

	if entity.NotMyMethod(err) {
		// send to proxy
		return ex.SendProxy(pkt)
	}

	return err
}
