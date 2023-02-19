package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
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
		// send to client
		id := pkt.Header.GetDstEntity()
		rw, ok := cltMap[entity.EntityID(id.ID)]
		if !ok {
			return errors.New("unknown client")
		}

		return rw.WritePacket(pkt)
	}

	return err

}
