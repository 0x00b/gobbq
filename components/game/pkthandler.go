package game

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
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
	return st.MethodPacketHandler.HandlePacket(pkt)
}
