package main

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/proto/bbq"
)

var _ nets.PacketHandler = &Proxy{}

func (p *Proxy) isMyPacket(pkt *nets.Packet) bool {

	hdr := pkt.Header
	dstEty := entity.DstEntity(pkt)
	if hdr.GetServiceType() == bbq.ServiceType_Entity ||
		hdr.RequestType == bbq.RequestType_RequestRespone {

		if dstEty.ProxyID() == p.EntityID().ProxyID() &&
			dstEty.InstID() == p.EntityID().InstID() {
			return true
		}
		return false
	}

	_, ok := p.EntityMgr.GetService(hdr.GetType())
	return ok
}

func (p *Proxy) HandlePacket(pkt *nets.Packet) error {

	// xlog.Debugln("proxy 1")

	hdr := pkt.Header
	dstEty := entity.DstEntity(pkt)

	if hdr.GetServiceType() == bbq.ServiceType_Entity && dstEty.Invalid() {
		// xlog.Errorln("bad req header:", hdr.String())
		return erro.ErrBadCall
	}

	if p.isMyPacket(pkt) {
		return p.Server.EntityMgr.HandlePacket(pkt)
	}

	// xlog.Debugln("proxy 2")

	// send to game or gate
	switch hdr.ServiceType {
	case bbq.ServiceType_Entity:
		return p.ProxyToEntity(dstEty, pkt)
	case bbq.ServiceType_Service:
		return p.ProxyToService(hdr, pkt)
	default:
		return erro.ErrNoServiveType
	}
}
