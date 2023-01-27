package main

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
)

type clientMap map[entity.EntityID]*codec.PacketReadWriter

var cltMap clientMap = make(clientMap)

// // RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	cltMap[sid] = prw
}

// func Login(pkt *codec.Packet) {
// 	// login
// 	// get entity id
// 	// register proxy
// 	id := snowflake.GenUUID()

// 	RegisterEntity(entity.EntityID(id), pkt.Src)

// 	RegisterProxy(entity.EntityID(id))

// }
