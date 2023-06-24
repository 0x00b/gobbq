package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &ClientPacketHandler{}

type ClientPacketHandler struct {
	gate *Gate
}

func (gt *Gate) isMyPacketFromClient(pkt *nets.Packet) bool {

	hdr := pkt.Header
	if hdr.RequestType == bbq.RequestType_RequestRequest &&
		hdr.GetServiceType() == bbq.ServiceType_Entity {
		return false
	}

	// 只会向客户端提供service， 所以不需要判断

	_, ok := gt.EntityMgr.GetService(hdr.GetType())
	return ok

}

func NewClientPacketHandler(gate *Gate) *ClientPacketHandler {
	st := &ClientPacketHandler{
		gate: gate,
	}
	return st
}

func (st *ClientPacketHandler) HandlePacket(pkt *nets.Packet) error {

	if st.gate.isMyPacketFromClient(pkt) {
		err := st.gate.EntityMgr.HandlePacket(pkt)
		if err != nil {
			xlog.Errorln("bad req handle:", pkt.Header.String(), err)
		}
		return err
	}

	// send to proxy
	return ex.SendProxy(pkt)

}
