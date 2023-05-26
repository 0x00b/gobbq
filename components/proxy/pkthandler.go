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

func (p *Proxy) HandlePacket(pkt *codec.Packet) error {

	xlog.Debugln("proxy 1")

	hdr := pkt.Header

	// todo
	// if p.EntityID().ID == pkt.Header.DstEntity.InstID

	err := p.Server.EntityMgr.HandlePacket(pkt)
	if err == nil {
		// handle succ
		return nil
	}

	xlog.Debugln("proxy 2")

	if entity.NotMyMethod(err) {
		// request
		// send to game
		// or send to gate
		if hdr.ServiceType == bbq.ServiceType_Entity {
			if hdr.DstEntity == nil {
				xlog.Errorln("bad req header:", hdr.String())
				return errors.New("bad call, call entity but no dst entity")
			}
			xlog.Debugln("proxy 3")
			p.ProxyToEntity(hdr.DstEntity, pkt)
			xlog.Debugln("proxy 4")
		} else {
			// call service
			p.ProxyToService(hdr, pkt)
		}

		return nil
	}

	return err
}
