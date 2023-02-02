package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ClientPacketHandler{}

type ClientPacketHandler struct {
	*entity.MethodPacketHandler
}

func NewClientPacketHandler() *ClientPacketHandler {
	st := &ClientPacketHandler{entity.NewMethodPacketHandler()}
	return st
}

func (st *ClientPacketHandler) HandlePacket(pkt *codec.Packet) error {

	err := st.MethodPacketHandler.HandlePacket(pkt)
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
