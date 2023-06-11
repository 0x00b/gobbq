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

	p.EntityMgr.RegisterEntity(nil, entity.FixedEntityID(entity.ProxyID(eid), entity.InstID(eid), entity.ID(eid)), p)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go p.Run(&wg)
	wg.Wait()

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

type ProxyMap map[entity.ProxyID]*codec.PacketReadWriter
type ProxySvcMap map[string]*codec.PacketReadWriter
type instMap map[entity.InstID]*codec.PacketReadWriter
type serviceMap map[string]*codec.PacketReadWriter

// // RegisterEntity register serive
func (p *Proxy) registerInst(instID entity.InstID, prw *codec.PacketReadWriter) {
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

func (p *Proxy) ProxyToEntity(eid entity.EntityID, pkt *codec.Packet) {

	// xlog.Debugln("proxy 11")
	// proxy to inst
	sendInst := func() bool {
		xlog.Debugln("local proxy to id:", eid)
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

func (p *Proxy) ProxyToService(hdr *bbq.Header, pkt *codec.Packet) {

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

		entity.SetRemoteEntityManager(c, prxy)

		xlog.Println("RegisterProxy ing...")

		rsp, err := proxypb.NewProxyEtyEntityClient(entity.FixedEntityID(entity.ProxyID(cfg.ID), entity.InstID(cfg.ID), entity.ID(cfg.ID))).
			RegisterProxy(c, &proxypb.RegisterProxyRequest{
				ProxyID: uint32(p.EntityID().ID()),
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
