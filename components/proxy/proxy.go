package main

import (
	"sync"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

func NewProxy() *Proxy {

	conf.Init("proxy.yaml")

	p := &Proxy{
		etyMtx:      sync.RWMutex{},
		entityMaps:  make(entityMap),
		svcMtx:      sync.RWMutex{},
		svcMaps:     make(serviceMap),
		proxyMap:    make(ProxyMap),
		proxySvcMap: make(ProxySvcMap),

		Server: bs.NewServer(),
	}

	p.Server.EntityMgr.EntityIDGenerator = p

	desc := proxypb.ProxyEtyEntityDesc
	desc.EntityImpl = p
	desc.EntityMgr = p.Server.EntityMgr
	p.SetDesc(&desc)

	eid := &bbq.EntityID{ID: conf.C.Proxy.Inst[0].ID, Type: proxypb.ProxyEtyEntityDesc.TypeName}

	p.Server.EntityMgr.RegisterEntity(nil, eid, p)

	go p.Run()

	p.ConnOtherProxy(nets.WithPacketHandler(p))

	return p
}

type Proxy struct {
	entity.Entity

	etyMtx     sync.RWMutex
	entityMaps entityMap

	svcMtx  sync.RWMutex
	svcMaps serviceMap

	proxyMap ProxyMap

	proxySvcMtx sync.RWMutex
	proxySvcMap ProxySvcMap

	Server *bs.Server
}

type ProxyMap map[string]*codec.PacketReadWriter
type ProxySvcMap map[string]*codec.PacketReadWriter
type entityMap map[string]*codec.PacketReadWriter
type serviceMap map[string]*codec.PacketReadWriter

// RegisterEntity register serive
func (p *Proxy) registerEntity(eid *bbq.EntityID, prw *codec.PacketReadWriter) {
	p.etyMtx.Lock()
	defer p.etyMtx.Unlock()
	if _, ok := p.entityMaps[eid.ID]; ok {
		xlog.Errorln("already has entity", eid)
	}
	xlog.Debugln("register entity id:", eid)
	p.entityMaps[eid.ID] = prw
}

// RegisterEntity register serive
func (p *Proxy) registerService(svcName string, prw *codec.PacketReadWriter) {
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

	xlog.Debugln("proxy 11")
	// proxy to local
	sendLocal := func() bool {
		p.etyMtx.RLock()
		defer p.etyMtx.RUnlock()
		prw, ok := p.entityMaps[eid.ID]
		if ok {
			xlog.Debugln("proxy to id:", eid)
			prw.SendPackt(pkt)
			xlog.Debugln("proxy 22")
			return true
		}
		return false
	}()

	xlog.Debugln("proxy 33")
	if sendLocal {
		return
	}

	// proxy to other proxy
	sendProxy := func() bool {
		proxyID := pkt.Header.DstEntity.ProxyID
		prw, ok := p.proxyMap[proxyID]
		if ok {
			prw.SendPackt(pkt)
			xlog.Debugln("proxy 44")
			return true
		}
		xlog.Errorln("unknown proxyid", proxyID)
		return false
	}()

	xlog.Debugln("proxy 55")

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
			xlog.Debugln("proxy to svc", prw)
			prw.SendPackt(pkt)
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
			prw.SendPackt(pkt)
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

		c, release := p.Context().Copy()
		defer release()

		entity.SetRemoteEntityManager(c, prxy)

		rsp, err := proxypb.NewProxyEtyEntityClient(&bbq.EntityID{ID: cfg.ID, ProxyID: cfg.ID}).
			RegisterProxy(c, &proxypb.RegisterProxyRequest{
				ProxyID: string(p.EntityID().ID),
			})
		if err != nil {
			panic(err)
		}

		p.proxyMap[cfg.ID] = prxy.GetPacketReadWriter()
		for _, v := range rsp.SvcNames {
			p.RegisterProxyService(v, prxy.GetPacketReadWriter())
		}
	}
}

func (p *Proxy) NewEntityID(tn string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: tn, ProxyID: p.EntityID().ID}
}
