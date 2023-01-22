package main

import (
	"github.com/0x00b/gobbq/components/proxy/px"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
)

type clientMap map[entity.EntityID]*codec.PacketReadWriter

type clientProxy map[entity.EntityID]entity.NopEntity

var cltmap clientMap

var proxymap px.ProxyMap

// RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	cltmap[sid] = prw
}

func Recv(sid entity.EntityID, pkt *codec.Packet) {

	prw := proxymap[sid]

	prw.WritePacket(pkt)
}

func Send(sid entity.EntityID, pkt *codec.Packet) {

	prw := cltmap[sid]

	prw.WritePacket(pkt)
}

func Login() {
	// login
	// get entity id
	// register proxy
}
