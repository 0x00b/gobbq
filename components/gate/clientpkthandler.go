package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &ClientPacketHandler{}

type ClientPacketHandler struct {
	gate *Gate
}

func NewClientPacketHandler(gate *Gate) *ClientPacketHandler {
	st := &ClientPacketHandler{
		gate: gate,
	}
	return st
}

func (st *ClientPacketHandler) HandlePacket(pkt *nets.Packet) error {

	if st.gate.isMyPacket(pkt) {
		return st.gate.EntityMgr.HandlePacket(pkt)
	}

	// send to proxy
	return ex.SendProxy(pkt)

}
