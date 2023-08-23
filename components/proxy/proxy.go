package main

import (
	"context"
	"sync"
	"sync/atomic"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/db/mongo"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
	capi "github.com/hashicorp/consul/api"
)

func NewProxy() *Proxy {

	p := &Proxy{
		instMtx:     sync.RWMutex{},
		instMaps:    make(instMap),
		svcMaps:     make(serviceMap),
		proxyMap:    make(ProxyMap),
		proxySvcMap: make(ProxySvcMap),

		Server: bs.NewServer(),
	}

	p.EntityMgr.EntityIDGenerator = p

	desc := proxypb.ProxyEtyEntityDesc
	desc.EntityImpl = p
	desc.EntityMgr = p.EntityMgr
	entity.SetEntityDesc(p, &desc)

	p.db = mongo.NewMongoDB()

	err := p.db.Connect(CFG.Mongo)
	if err != nil {
		panic(err)
	}

	eid, err := p.db.GetIncrementID(context.Background(), "gobbq_proxy_inst_id")
	if err != nil {
		panic(err)
	}

	xlog.Infoln("proxy entity id", eid)

	proxypb.RegisterProxyEtyEntity(p.EntityMgr, p)
	p.EntityMgr.RegisterEntity(nil, entity.FixedEntityID(entity.ProxyID(eid), 0, 0), p)

	entity.Run(p)

	p.consul, err = capi.NewClient(capi.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// 连接其他proxy
	p.ConnOtherProxys(nets.WithPacketHandler(p), nets.WithConnCallback(p))

	// 注册自己到consul
	p.RegisterToConsul()

	return p
}

type Proxy struct {
	*bs.Server

	entity.Entity

	// 属于service，需要加锁
	instMtx  sync.RWMutex
	instMaps instMap

	// 下面的都是entity使用，不需要锁
	// svcMtx  sync.RWMutex
	svcMaps serviceMap
	// proxyMtx sync.RWMutex
	proxyMap ProxyMap
	// proxySvcMtx sync.RWMutex
	proxySvcMap ProxySvcMap

	instIdCounter uint32

	db db.IDatabase

	consul *capi.Client
}

type ProxyMap map[entity.ProxyID]*nets.Conn
type ProxySvcMap map[string][]*nets.Conn
type instMap map[entity.InstID]*nets.Conn
type serviceMap map[string][]*nets.Conn

func (p *Proxy) NewInstID() entity.InstID {
	id := atomic.AddUint32(&p.instIdCounter, 1)

	if id == 0 {
		// report
		panic("cycle new inst id")
	}

	return entity.InstID(id)
}

// RegisterEntity register serive
func (p *Proxy) registerInst(instID entity.InstID, prw *nets.Conn) {
	p.instMtx.Lock()
	defer p.instMtx.Unlock()
	if _, ok := p.instMaps[instID]; ok {
		xlog.Traceln("already has entity", instID)
	}
	xlog.Traceln("register entity id:", instID)
	p.instMaps[instID] = prw
}

// getInst get inst
func (p *Proxy) getInst(instID entity.InstID) (*nets.Conn, bool) {
	p.instMtx.RLock()
	defer p.instMtx.RUnlock()

	cn, ok := p.instMaps[instID]
	return cn, ok
}

// RegisterEntity register serive
func (p *Proxy) registerService(svcName string, prw *nets.Conn) {
	// p.svcMtx.Lock()
	// defer p.svcMtx.Unlock()

	p.svcMaps[svcName] = append(p.svcMaps[svcName], prw)
}

// RegisterEntity register serive
func (p *Proxy) getService(svcName string) ([]*nets.Conn, bool) {
	// p.svcMtx.RLock()
	// defer p.svcMtx.RUnlock()

	cn, ok := p.svcMaps[svcName]
	return cn, ok
}

// RegisterEntity register serive
func (p *Proxy) RegisterProxyService(svcName string, prw *nets.Conn) {
	// p.proxySvcMtx.Lock()
	// defer p.proxySvcMtx.Unlock()
	if _, ok := p.proxySvcMap[svcName]; ok {
		xlog.Traceln("already has svc")
	}

	p.proxySvcMap[svcName] = append(p.proxySvcMap[svcName], prw)
}

// RegisterEntity register serive
func (p *Proxy) getProxyService(svcName string) ([]*nets.Conn, bool) {
	// p.proxySvcMtx.RLock()
	// defer p.proxySvcMtx.RUnlock()

	cn, ok := p.proxySvcMap[svcName]
	return cn, ok
}

func (p *Proxy) ProxyToEntity(eid entity.EntityID, pkt *nets.Packet) (err error) {

	// proxy to inst
	sendInst := func() bool {
		if eid.ProxyID() != p.EntityID().ProxyID() {
			return false
		}

		prw, ok := p.getInst(eid.InstID())
		if ok {
			err = prw.SendPacket(pkt)
		} else {
			err = erro.ErrProxyInstNotFound
		}
		return true
	}()

	if sendInst {
		return
	}

	// proxy to other proxy
	func() bool {
		proxyID := entity.DstEntity(pkt).ProxyID()
		// p.proxyMtx.RLock()
		// defer p.proxyMtx.RUnlock()
		prw, ok := p.proxyMap[proxyID]
		if ok {
			err = prw.SendPacket(pkt)
			return true
		}
		err = erro.ErrProxyNotFound
		return false
	}()

	return err
}

func (p *Proxy) ProxyToService(hdr *bbq.Header, pkt *nets.Packet) (err error) {

	if hdr.RequestType == bbq.RequestType_RequestRespone {
		return p.ProxyToEntity(entity.DstEntity(pkt), pkt)
	}

	svcType := hdr.Type
	// proxy to local
	sendLocal := func() bool {

		prws, ok := p.getService(svcType) // p.svcMaps[svcType]
		if !ok || len(prws) == 0 {
			return false
		}
		sid := hdr.GetSrcEntity()
		err = prws[sid%uint64(len(prws))].SendPacket(pkt)
		return true
	}()

	if sendLocal {
		return
	}

	// proxy to other proxy
	func() bool {
		typ := pkt.Header.Type
		prws, ok := p.getProxyService(typ)
		if !ok || len(prws) == 0 {
			err = erro.ErrProxyServiceNotFound
			return false
		}
		sid := hdr.GetSrcEntity()
		err = prws[sid%uint64(len(prws))].SendPacket(pkt)
		return true
	}()

	return
}

func (p *Proxy) NewEntityID() entity.EntityID {
	return entity.NewEntityID(p.EntityID().ProxyID(), p.EntityID().InstID())
}

func (p *Proxy) Serve() error {
	return p.Server.ListenAndServe()
}
