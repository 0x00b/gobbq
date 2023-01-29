package main

import (
	"context"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
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

type ProxyService struct {
	entity.Service
}

// RegisterEntity
func (ps *ProxyService) RegisterEntity(c context.Context, req *proxypb.RegisterEntityRequest, ret func(*proxypb.RegisterEntityResponse, error)) {

	RegisterEntity(entity.EntityID(req.EntityID), nil)

	// //////

	cli := proxypb.NewProxyService()

	err := cli.RegisterEntity(c, req, func(c context.Context, rsp *proxypb.RegisterEntityResponse) {

		ret(rsp, nil)

	})

	if err != nil {
		ret(nil, err)
	}

	return
}

// UnregisterEntity
func (ps *ProxyService) UnregisterEntity(c context.Context, req *proxypb.RegisterEntityRequest, ret func(*proxypb.RegisterEntityResponse, error)) {
	return
}

// Ping
func (ps *ProxyService) Ping(c context.Context, req *proxypb.PingPong, ret func(*proxypb.PingPong, error)) {
	return
}
