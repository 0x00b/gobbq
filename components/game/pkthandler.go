package game

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &GamePacketHandler{}

type GamePacketHandler struct {
	*entity.MethodPacketHandler
}

func NewGamePacketHandler() *GamePacketHandler {
	st := &GamePacketHandler{entity.NewMethodPacketHandler()}
	return st
}

func (st *GamePacketHandler) HandlePacket(pkt *codec.Packet) error {
	hdr := pkt.Header
	if hdr.RequestType == bbq.RequestType_RequestRequest {
		return entity.NewMethodPacketHandler().HandlePacket(pkt)
	}
	// response
	xlog.Println("recv response:", pkt.String())

	return nil
}
