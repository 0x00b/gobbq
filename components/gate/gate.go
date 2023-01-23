package main

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
)

type clientMap map[entity.EntityID]*codec.PacketReadWriter

var cltmap clientMap

var proxymap ex.ProxyMap

// RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	cltmap[sid] = prw
}

// RegisterEntity register serive
func RegisterProxy(sid entity.EntityID) {

}

func Recv(sid entity.EntityID, pkt *codec.Packet) {

	prw := proxymap[sid]

	prw.WritePacket(pkt)
}

func Send(sid entity.EntityID, pkt *codec.Packet) {

	prw := cltmap[sid]

	prw.WritePacket(pkt)
}

func Login(pkt *codec.Packet) {
	// login
	// get entity id
	// register proxy
	id := snowflake.GenUUID()

	RegisterEntity(entity.EntityID(id), pkt.Src)

	RegisterProxy(entity.EntityID(id))

}
