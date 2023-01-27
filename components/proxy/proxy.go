package main

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
)

type entityMap map[entity.EntityID]*codec.PacketReadWriter

var entityMaps entityMap = make(entityMap)

// RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	entityMaps[sid] = prw
}

func Proxy(sid entity.EntityID, pkt *codec.Packet) {

	prw := entityMaps[sid]

	prw.WritePacket(pkt)
}
