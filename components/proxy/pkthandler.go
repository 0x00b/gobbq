package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &ProxyPacketHandler{}

type ProxyPacketHandler struct {
	*entity.MethodPacketHandler
}

func NewProxyPacketHandler() *ProxyPacketHandler {
	st := &ProxyPacketHandler{entity.NewMethodPacketHandler()}
	return st
}

func (st *ProxyPacketHandler) HandlePacket(pkt *codec.Packet) error {
	hdr := pkt.Header
	err := st.MethodPacketHandler.HandlePacket(pkt)
	if err == nil {
		// handle succ
		return nil
	}

	if entity.NotMyMethod(err) {
		// request
		// send to game
		// or send to gate
		if hdr.ServiceType == bbq.ServiceType_Entity {
			if hdr.DstEntity == nil {
				xlog.Errorln("bad req header:", hdr.String())
				return errors.New("bad call, call entity but no dst entity")
			}
			proxyInst.ProxyToEntity(hdr.DstEntity, pkt)
		} else {
			// call service
			proxyInst.ProxyToService(hdr, pkt)
		}

		return nil
	}

	return err
}
