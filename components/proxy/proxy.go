package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
)

type entityMap map[entity.EntityID]*codec.PacketReadWriter
type serviceMap map[string][]*codec.PacketReadWriter

var entityMaps entityMap = make(entityMap)
var svcMap serviceMap = make(serviceMap)

// RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	entityMaps[sid] = prw
}

// RegisterEntity register serive
func RegisterService(svcName string, prw *codec.PacketReadWriter) {
	svcMap[svcName] = append(svcMap[svcName], prw)
}

func ProxyToEntity(sid entity.EntityID, pkt *codec.Packet) {
	prw, ok := entityMaps[sid]
	if !ok {
		fmt.Println("unknown entity id")
		return
	}
	prw.WritePacket(pkt)
}

func ProxyToService(hdr *bbq.Header, pkt *codec.Packet) {
	prws, ok := svcMap[hdr.ServiceName]
	if !ok {
		fmt.Println("unknown entity id")
		return
	}
	if hdr.SrcEntity != nil {
		// hash
		prws[0].WritePacket(pkt)
	}
	// random
	prws[0].WritePacket(pkt)
}

type ProxyService struct {
	entity.Service
}

// RegisterEntity
func (ps *ProxyService) RegisterEntity(c *entity.Context, req *proxypb.RegisterEntityRequest, ret func(*proxypb.RegisterEntityResponse, error)) {

	fmt.Println("register entity id:", req.EntityID)
	RegisterEntity(entity.EntityID(req.EntityID), c.Packet().Src)

	return
}

// RegisterEntity
func (ps *ProxyService) RegisterService(c *entity.Context, req *proxypb.RegisterServiceRequest, ret func(*proxypb.RegisterServiceResponse, error)) {

	fmt.Println("register service:", req.ServiceName)
	RegisterService(req.ServiceName, c.Packet().Src)

	return
}

// UnregisterEntity
func (ps *ProxyService) UnregisterEntity(c *entity.Context, req *proxypb.RegisterEntityRequest, ret func(*proxypb.RegisterEntityResponse, error)) {
	return
}

// RegisterEntity
func (ps *ProxyService) UnregisterService(c *entity.Context, req *proxypb.RegisterServiceRequest, ret func(*proxypb.RegisterServiceResponse, error)) {
	return
}

// Ping
func (ps *ProxyService) Ping(c *entity.Context, req *proxypb.PingPong, ret func(*proxypb.PingPong, error)) {
	return
}

type ProxyMap map[uint32]*nets.Client

var proxyMap ProxyMap = make(ProxyMap)

func ConnProxy(ops ...nets.Option) {
	for i := 0; i < int(conf.C.Proxy.InstNum); i++ {

		// connect to proxy
		cfg := conf.C.Proxy.Inst[i]
		_ = cfg.Net

		// jump myself
		// if cfg.ID == 0 {

		// }

		prxy, err := nets.Connect(nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, ops...)

		if err != nil {
			panic(err)
		}

		proxyMap[cfg.ID] = prxy
	}
}

func SendProxy(pkt *codec.Packet) error {
	_ = pkt.Header.GetDstEntity().ID
	// hash id , lb proxy

	return proxyMap[1].SendPackt(pkt)

}
