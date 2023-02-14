package main

import (
	"sync"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

type entityMap map[entity.EntityID]*codec.PacketReadWriter
type serviceMap map[string][]*codec.PacketReadWriter

var entityMaps entityMap = make(entityMap)
var svcMap serviceMap = make(serviceMap)
var etyMtx sync.RWMutex

// RegisterEntity register serive
func RegisterEntity(eid entity.EntityID, prw *codec.PacketReadWriter) {
	etyMtx.Lock()
	defer etyMtx.Unlock()
	xlog.Println("register entity id:", eid)
	entityMaps[eid] = prw
}

// RegisterEntity register serive
func RegisterService(svcName string, prw *codec.PacketReadWriter) {
	etyMtx.RLock()
	defer etyMtx.RUnlock()
	svcMap[svcName] = append(svcMap[svcName], prw)
}

func ProxyToEntity(eid entity.EntityID, pkt *codec.Packet) {
	etyMtx.RLock()
	defer etyMtx.RUnlock()
	prw, ok := entityMaps[eid]
	if !ok {
		xlog.Println("unknown entity id", eid)
		return
	}
	xlog.Println("proxy to id:", eid)
	prw.WritePacket(pkt)
}

func ProxyToService(hdr *bbq.Header, pkt *codec.Packet) {
	etyMtx.RLock()
	defer etyMtx.RUnlock()
	prws, ok := svcMap[hdr.ServiceName]
	if !ok {
		xlog.Println("unknown svc name", hdr.ServiceName)
		return
	}
	xlog.Println("proxy to svc", len(prws), prws)
	if hdr.SrcEntity != "" {
		// hash
		prws[0].WritePacket(pkt)
		return
	}
	// random
	prws[0].WritePacket(pkt)
}

type ProxyService struct {
	entity.Service
}

func (ps *ProxyService) OnInit() {
	xlog.Println("on init ProxyService")
}

// RegisterInst
func (ps *ProxyService) RegisterInst(c entity.Context, req *proxypb.RegisterInstRequest) (*proxypb.RegisterInstResponse, error) {

	RegisterEntity(entity.EntityID(req.EntityID), c.Packet().Src)

	return &proxypb.RegisterInstResponse{}, nil
}

// RegisterEntity
func (ps *ProxyService) RegisterEntity(c entity.Context, req *proxypb.RegisterEntityRequest) (*proxypb.RegisterEntityResponse, error) {

	RegisterEntity(entity.EntityID(req.EntityID), c.Packet().Src)

	return &proxypb.RegisterEntityResponse{}, nil
}

// RegisterEntity
func (ps *ProxyService) RegisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	xlog.Println("register service:", req.ServiceName)
	RegisterService(req.ServiceName, c.Packet().Src)

	return &proxypb.RegisterServiceResponse{}, nil
}

// UnregisterEntity
func (ps *ProxyService) UnregisterEntity(c entity.Context, req *proxypb.RegisterEntityRequest) (*proxypb.RegisterEntityResponse, error) {

	return nil, nil
}

// RegisterEntity
func (ps *ProxyService) UnregisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	return nil, nil
}

// Ping
func (ps *ProxyService) Ping(c entity.Context, req *proxypb.PingPong) (*proxypb.PingPong, error) {

	return nil, nil
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
	_ = pkt.Header.GetDstEntity()
	// hash id , lb proxy

	return proxyMap[1].SendPackt(pkt)

}
