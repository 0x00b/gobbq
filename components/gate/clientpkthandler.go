package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
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
		err := st.gate.EntityMgr.HandlePacket(pkt)
		if err != nil {
			xlog.Errorln("bad req handle:", pkt.Header.String(), err)
		}
		return err
	}

	// send to proxy
	return ex.SendProxy(pkt)

}
