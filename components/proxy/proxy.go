package main

import (
	"sync"
	"sync/atomic"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
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
	p.SetEntityDesc(&desc)

	eid := uint16(conf.C.Proxy.Inst[0].ID)

	proxypb.RegisterProxyEtyEntity(p.EntityMgr, p)
	p.EntityMgr.RegisterEntity(nil, entity.FixedEntityID(entity.ProxyID(eid), 0, 0), p)

	p.Run()

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

	instIdCounter uint32
}

type ProxyMap map[entity.ProxyID]*nets.PacketReadWriter
type ProxySvcMap map[string]*nets.PacketReadWriter
type instMap map[entity.InstID]*nets.PacketReadWriter
type serviceMap map[string]*nets.PacketReadWriter

func (p *Proxy) NewInstID() entity.InstID {
	id := atomic.AddUint32(&p.instIdCounter, 1)

	if id == 0 {
		// report
		panic("cycle new inst id")
	}

	return entity.InstID(id)
}

// // RegisterEntity register serive
func (p *Proxy) registerInst(instID entity.InstID, prw *nets.PacketReadWriter) {
	p.instMtx.Lock()
	defer p.instMtx.Unlock()
	if _, ok := p.instMaps[instID]; ok {
		xlog.Errorln("already has entity", instID)
	}
	xlog.Debugln("register entity id:", instID)
	p.instMaps[instID] = prw
}

// RegisterEntity register serive
func (p *Proxy) registerService(svcName string, prw *nets.PacketReadWriter) {
	p.svcMtx.RLock()
	defer p.svcMtx.RUnlock()
	if _, ok := p.svcMaps[svcName]; ok {
		xlog.Errorln("already has svc")
	}

	p.svcMaps[svcName] = prw
}

// RegisterEntity register serive
func (p *Proxy) RegisterProxyService(svcName string, prw *nets.PacketReadWriter) {
	p.proxySvcMtx.Lock()
	defer p.proxySvcMtx.Unlock()
	if _, ok := p.proxySvcMap[svcName]; ok {
		xlog.Errorln("already has svc")
	}

	p.proxySvcMap[svcName] = prw
}

func (p *Proxy) ProxyToEntity(eid entity.EntityID, pkt *nets.Packet) {

	// xlog.Debugln("proxy 11")
	// proxy to inst
	sendInst := func() bool {
		xlog.Debugln("local proxy to id:", eid)
		if eid.ProxyID() != p.EntityID().ProxyID() {
			xlog.Debugln("not my inst", eid)
			return false
		}
		p.instMtx.RLock()
		defer p.instMtx.RUnlock()
		prw, ok := p.instMaps[eid.InstID()]
		if ok {
			xlog.Debugln("proxy to id:", eid)
			prw.SendPacket(pkt)
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
		proxyID := entity.DstEntity(pkt).ProxyID()
		prw, ok := p.proxyMap[proxyID]
		if ok {
			prw.SendPacket(pkt)
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

func (p *Proxy) ProxyToService(hdr *bbq.Header, pkt *nets.Packet) {

	if hdr.RequestType == bbq.RequestType_RequestRespone {
		p.ProxyToEntity(entity.DstEntity(pkt), pkt)
		return
	}

	svcType := hdr.Type
	// proxy to local
	sendLocal := func() bool {
		p.svcMtx.RLock()
		defer p.svcMtx.RUnlock()

		prw, ok := p.svcMaps[svcType]
		if ok {
			xlog.Debugln("proxy to svc", prw)
			prw.SendPacket(pkt)
			return true
		}
		return false
	}()

	if sendLocal {
		return
	}

	// proxy to other proxy
	sendProxy := func() bool {
		typ := pkt.Header.Type
		prw, ok := p.proxySvcMap[typ]
		if ok {
			prw.SendPacket(pkt)
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

		entity.SetProxy(c, prxy)

		xlog.Println("RegisterProxy ing...")

		rsp, err := proxypb.NewProxyEtyEntityClient(entity.FixedEntityID(entity.ProxyID(cfg.ID), 0, 0)).
			RegisterProxy(c, &proxypb.RegisterProxyRequest{
				ProxyID: uint32(p.EntityID().ProxyID()),
			})
		if err != nil {
			panic(err)
		}
		xlog.Println("RegisterProxy done...")

		p.proxyMap[entity.ProxyID(cfg.ID)] = prxy.GetPacketReadWriter()
		for _, v := range rsp.ServiceNames {
			p.RegisterProxyService(v, prxy.GetPacketReadWriter())
		}
	}
}

func (p *Proxy) NewEntityID() entity.EntityID {
	return entity.NewEntityID(p.EntityID().ProxyID(), p.EntityID().InstID())
}

func (p *Proxy) Serve() error {
	return p.Server.ListenAndServe()
}
