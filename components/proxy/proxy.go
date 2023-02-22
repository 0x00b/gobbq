package main

import (
	"sync"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

type Proxy struct {
	entity.Entity

	etyMtx     sync.RWMutex
	entityMaps entityMap

	svcMtx  sync.RWMutex
	svcMaps serviceMap

	proxyMap ProxyMap

	proxySvcMtx sync.RWMutex
	proxySvcMap ProxySvcMap
}

var proxyInst = NewProxy()

func NewProxy() *Proxy {

	entity.NewEntityID = &EntityIDGenerator{}

	gm := &Proxy{
		etyMtx:      sync.RWMutex{},
		entityMaps:  make(entityMap),
		svcMtx:      sync.RWMutex{},
		svcMaps:     make(serviceMap),
		proxyMap:    make(ProxyMap),
		proxySvcMap: make(ProxySvcMap),
	}
	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, &bbq.EntityID{ID: eid}, gm)

	go gm.Run()

	return gm
}

type ProxyMap map[string]*codec.PacketReadWriter
type ProxySvcMap map[string]*codec.PacketReadWriter
type entityMap map[string]*codec.PacketReadWriter
type serviceMap map[string]*codec.PacketReadWriter

// RegisterEntity register serive
func (p *Proxy) RegisterEntity(eid *bbq.EntityID, prw *codec.PacketReadWriter) {
	p.etyMtx.Lock()
	defer p.etyMtx.Unlock()
	if _, ok := p.entityMaps[eid.ID]; ok {
		xlog.Errorln("already has entity", eid)
	}
	xlog.Println("register entity id:", eid)
	p.entityMaps[eid.ID] = prw
}

// RegisterEntity register serive
func (p *Proxy) RegisterService(svcName string, prw *codec.PacketReadWriter) {
	p.svcMtx.RLock()
	defer p.svcMtx.RUnlock()
	if _, ok := p.svcMaps[svcName]; ok {
		xlog.Errorln("already has svc")
	}

	p.svcMaps[svcName] = prw
}

// RegisterEntity register serive
func (p *Proxy) RegisterProxyService(svcName string, prw *codec.PacketReadWriter) {
	p.proxySvcMtx.Lock()
	defer p.proxySvcMtx.Unlock()
	if _, ok := p.proxySvcMap[svcName]; ok {
		xlog.Errorln("already has svc")
	}

	p.proxySvcMap[svcName] = prw
}

func (p *Proxy) ProxyToEntity(eid *bbq.EntityID, pkt *codec.Packet) {

	// proxy to local
	sendLocal := func() bool {
		p.etyMtx.RLock()
		defer p.etyMtx.RUnlock()
		prw, ok := p.entityMaps[eid.ID]
		if ok {
			xlog.Println("proxy to id:", eid)
			prw.WritePacket(pkt)
			return true
		}
		return false
	}()

	if sendLocal {
		return
	}

	// proxy to other proxy
	sendProxy := func() bool {
		proxyID := pkt.Header.DstEntity.ProxyID
		prw, ok := p.proxyMap[proxyID]
		if ok {
			prw.WritePacket(pkt)
			return true
		}
		xlog.Errorln("unknown proxyid", proxyID)
		return false
	}()

	if !sendProxy {
		xlog.Errorln("unknown entity", eid)
	}
}

func (p *Proxy) ProxyToService(hdr *bbq.Header, pkt *codec.Packet) {

	if hdr.RequestType == bbq.RequestType_RequestRespone {
		p.ProxyToEntity(pkt.Header.DstEntity, pkt)
		return
	}

	svcType := hdr.DstEntity.Type
	// proxy to local
	sendLocal := func() bool {
		p.svcMtx.RLock()
		defer p.svcMtx.RUnlock()

		prw, ok := p.svcMaps[svcType]
		if ok {
			xlog.Println("proxy to svc", prw)
			prw.WritePacket(pkt)
			return true
		}
		return false
	}()

	if sendLocal {
		return
	}

	// proxy to other proxy
	sendProxy := func() bool {
		typ := pkt.Header.DstEntity.Type
		prw, ok := p.proxySvcMap[typ]
		if ok {
			prw.WritePacket(pkt)
			return true
		}
		return false
	}()

	if !sendProxy {
		xlog.Errorln("unknown svc", svcType)
	}
}

func (p *Proxy) ConnOtherProxy(ops ...nets.Option) {
	for i := 1; i < int(conf.C.Proxy.InstNum); i++ {

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

		rsp, err := proxypb.NewProxyServiceClient(prxy.GetPacketReadWriter()).RegisterProxy(proxyInst.Context(), &proxypb.RegisterProxyRequest{
			ProxyID: string(proxyInst.EntityID().ID),
		})
		if err != nil {
			panic(err)
		}

		p.proxyMap[rsp.ProxyID] = prxy.GetPacketReadWriter()
		for _, v := range rsp.SvcNames {
			p.RegisterProxyService(v, prxy.GetPacketReadWriter())
		}

	}
}

type EntityIDGenerator struct {
}

func (n *EntityIDGenerator) NewEntityID(tn string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: tn, ProxyID: proxyInst.EntityID().ID}
}
