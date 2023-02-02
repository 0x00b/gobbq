package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
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

	err := st.MethodPacketHandler.HandlePacket(pkt)
	if err == nil {
		// handle succ
		return nil
	}

	if entity.NotMyMethod(err) {
		hdr := pkt.Header
		// send to game
		// or send to gate
		if hdr.ServiceType == bbq.ServiceType_Entity {
			if hdr.DstEntity == nil {
				return errors.New("bad call, call entity but no dst entity")
			}
			ProxyToEntity(entity.EntityID(hdr.DstEntity.ID), pkt)
			return nil
		}
		// call service
		ProxyToService(hdr, pkt)
	}

	return err
}
