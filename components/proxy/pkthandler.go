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
	var err error
	hdr := pkt.Header
	if hdr.RequestType == bbq.RequestType_RequestRequest {
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
				if hdr.DstEntity == "" {
					xlog.Println("bad req header:", hdr.String())
					return errors.New("bad call, call entity but no dst entity")
				}
				ProxyToEntity(entity.EntityID(hdr.DstEntity), pkt)
			} else {
				// call service
				ProxyToService(hdr, pkt)
			}

			return nil
		}
	}
	// response
	ProxyToEntity(entity.EntityID(hdr.DstEntity), pkt)

	return err
}
