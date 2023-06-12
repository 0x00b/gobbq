package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &Proxy{}

func (p *Proxy) isMyPacket(pkt *codec.Packet) bool {

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

	return p.EntityMgr.IsMyService(hdr.GetType())
}

func (p *Proxy) HandlePacket(pkt *codec.Packet) error {

	// xlog.Debugln("proxy 1")

	hdr := pkt.Header
	dstEty := entity.DstEntity(pkt)

	if hdr.GetServiceType() == bbq.ServiceType_Entity && dstEty.Invalid() {
		xlog.Errorln("bad req header:", hdr.String())
		return errors.New("bad call, call entity but no dst entity")
	}

	if p.isMyPacket(pkt) {
		err := p.Server.EntityMgr.HandlePacket(pkt)
		if err != nil {
			xlog.Errorln("bad req handle:", hdr.String(), err)
		}
		return err
	}

	// xlog.Debugln("proxy 2")

	// request
	// send to game
	// or send to gate
	if hdr.ServiceType == bbq.ServiceType_Entity {
		// xlog.Debugln("proxy 3")
		p.ProxyToEntity(dstEty, pkt)
		// xlog.Debugln("proxy 4")
	} else {
		// call service
		p.ProxyToService(hdr, pkt)
	}

	return nil
}
