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
		instMtx:     sync.RWMutex{},
		instMaps:    make(instMap),
		svcMtx:      sync.RWMutex{},
		svcMaps:     make(serviceMap),
		proxyMap:    make(ProxyMap),
		proxySvcMap: make(ProxySvcMap),

		Server: bs.NewServer(),
	}

	p.EntityMgr.EntityIDGenerator = p

	desc := proxypb.ProxyEtyEntityDesc
	desc.EntityImpl = p
	desc.EntityMgr = p.EntityMgr
	p.SetDesc(&desc)

	eid := &bbq.EntityID{ID: conf.C.Proxy.Inst[0].ID, Type: proxypb.ProxyEtyEntityDesc.TypeName}

	p.EntityMgr.RegisterEntity(nil, eid, p)

	go p.Run()

	p.ConnOtherProxy(nets.WithPacketHandler(p))

	return p
}

type Proxy struct {
	*bs.Server

	entity.Entity

	instMtx  sync.RWMutex
	instMaps instMap

	svcMtx  sync.RWMutex
	svcMaps serviceMap

	proxyMap ProxyMap

	proxySvcMtx sync.RWMutex
	proxySvcMap ProxySvcMap
}

type ProxyMap map[string]*codec.PacketReadWriter
type ProxySvcMap map[string]*codec.PacketReadWriter
type instMap map[string]*codec.PacketReadWriter
type serviceMap map[string]*codec.PacketReadWriter

// // RegisterEntity register serive
func (p *Proxy) registerInst(instID string, prw *codec.PacketReadWriter) {
	p.instMtx.Lock()
	defer p.instMtx.Unlock()
	if _, ok := p.instMaps[instID]; ok {
		xlog.Errorln("already has entity", instID)
	}
	xlog.Debugln("register entity id:", instID)
	p.instMaps[instID] = prw
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

	// xlog.Debugln("proxy 11")
	// proxy to inst
	sendInst := func() bool {
		xlog.Debugln("local proxy to id:", eid)
		p.instMtx.RLock()
		defer p.instMtx.RUnlock()
		prw, ok := p.instMaps[eid.GetInstID()]
		if ok {
			xlog.Debugln("proxy to id:", eid)
			prw.SendPackt(pkt)
			// xlog.Debugln("proxy 22")
			return true
		}
		return false
	}()

	// xlog.Debugln("proxy 33")
	if sendInst {
		return
	}

	// proxy to other proxy
	sendProxy := func() bool {
		proxyID := pkt.Header.DstEntity.ProxyID
		prw, ok := p.proxyMap[proxyID]
		if ok {
			prw.SendPackt(pkt)
			// xlog.Debugln("proxy 44")
			return true
		}
		xlog.Errorln("unknown proxyid", proxyID)
		return false
	}()

	// xlog.Debugln("proxy 55")

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
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: tn, ProxyID: p.EntityID().ID, InstID: p.EntityID().ID}
}

func (p *Proxy) Serve() error {
	return p.Server.ListenAndServe()
}
