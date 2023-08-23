package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
	capi "github.com/hashicorp/consul/api"
)

func NewGate() *Gate {

	InitConfig()

	gt := &Gate{
		Server:  bs.NewServer(),
		cltMap:  make(clientMap),
		watcher: make(map[entity.EntityID]map[entity.EntityID]bool),
	}

	gt.EntityMgr.ProxyRegister = gt
	gt.EntityMgr.EntityIDGenerator = gt

	desc := gatepb.GateServiceDesc
	desc.EntityImpl = gt
	desc.EntityMgr = gt.EntityMgr
	entity.SetServiceDesc(gt, &desc)

	temp := gt.NewEntityID()

	gt.EntityMgr.RegisterEntity(nil, temp, gt)

	entity.Run(gt)

	gt.init(temp)

	return gt
}

func (gt *Gate) init(old entity.EntityID) {

	cfg := CFG.Proxy
	if !CFG.LocalProxy {
		var err error
		gt.consul, err = capi.NewClient(capi.DefaultConfig())
		if err != nil {
			panic(err)
		}

		_, svcs, err := gt.consul.Agent().AgentHealthServiceByName("gobbqproxy")
		if err != nil {
			panic(err)
		}
		if len(svcs) <= 0 {
			panic("no svc")
		}

		// todo 负载均衡，现在先随机一个
		svc := svcs[rand.Int()%len(svcs)]
		cfg = conf.NetConf{
			Net:  string(nets.TCP),
			IP:   svc.Service.Address,
			Port: fmt.Sprint(svc.Service.Port),
		}

		xlog.Infoln("consul get proxy:", cfg)
	}

	ex.ConnProxy(cfg, nets.WithPacketHandler(gt), nets.WithConnCallback(&ProxyConnCallBack{gate: gt}))

	gt.EntityMgr.Proxy = ex.ProxyClient

	client := proxypb.NewProxySvcClient()

	rsp, err := client.RegisterInst(gt.Context(), &proxypb.RegisterInstRequest{})
	if err != nil {
		xlog.Errorln("error:", err)
		panic(err)
	}

	// 更新entityID
	proxyID := entity.ProxyID(rsp.GetProxyID())
	instID := entity.InstID(rsp.GetInstID())
	newid := entity.FixedEntityID(proxyID, instID, gt.Service.EntityID().ID())
	gt.EntityMgr.ReplaceEntityID(old, newid)

	gt.AddTimer(10*time.Second, func() bool {
		client.Ping(gt.Context(), &proxypb.PingPong{})
		return true
	})
}

type Gate struct {
	entity.Service

	cltMtx sync.Mutex
	cltMap clientMap

	*bs.Server

	watcher map[entity.EntityID]map[entity.EntityID]bool

	consul *capi.Client
}

type clientMap map[entity.EntityID]*nets.Conn

// // RegisterEntity register serive
func (gt *Gate) RegisterEntity(eid entity.EntityID, prw *nets.Conn) {
	gt.cltMtx.Lock()
	defer gt.cltMtx.Unlock()

	gt.cltMap[eid] = prw
}

// GateService
type GateService struct {
	entity.Service
}

// RegisterClient
func (gt *Gate) RegisterClient(c entity.Context, req *gatepb.RegisterClientRequest) (*gatepb.RegisterClientResponse, error) {

	id := gt.NewClientEntityID()

	gt.RegisterEntity(id, c.Packet().Src)

	// client := proxypb.NewProxySvcServiceClient()
	// rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: req.EntityID})
	// if err != nil {
	// 	return nil, err
	// }
	// xlog.Debugln("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{
		EntityID: uint64(id),
	}, nil
}

// UnregisterClient
func (gt *Gate) UnregisterClient(c entity.Context, req *gatepb.RegisterClientRequest) error {
	return nil
}

// Ping
func (gt *Gate) Ping(c entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

// func (gt *Gate) RegisterEntityToProxy(eid entity.EntityID) error {
// 	client := proxypb.NewProxySvcServiceClient()

// 	_, err := client.RegisterEntity(gt.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
// 	if err != nil {
// 		return err
// 	}

// 	xlog.Debugln("register proxy entity resp")

// 	return nil
// }

func (gt *Gate) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxyEtyClient(entity.FixedEntityID(gt.EntityID().ProxyID(), 0, 0))

	_, err := client.RegisterService(gt.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debugln("register proxy service resp")

	return nil
}

const (
	IdBitNum = entity.IDBitNum - 1
)

func (gt *Gate) NewEntityID() entity.EntityID {
	id := entity.GenID() | (1 << IdBitNum) //最高位是1代表是gate的entity， 为0代表client
	return entity.FixedEntityID(gt.EntityID().ProxyID(), gt.EntityID().InstID(), entity.ID(id))
}

func (gt *Gate) NewClientEntityID() entity.EntityID {
	id := entity.GenID() & (1<<IdBitNum - 1) //最高位是1代表是gate的entity， 为0代表client
	return entity.FixedEntityID(gt.EntityID().ProxyID(), gt.EntityID().InstID(), entity.ID(id))
}

func (gt *Gate) IsMyEntity(id entity.EntityID) bool {
	return id.ProxyID() == gt.EntityID().ProxyID() &&
		id.InstID() == gt.EntityID().InstID() &&
		(id.ID()>>IdBitNum) == 1
}

func (gt *Gate) Serve() error {
	return gt.Server.ListenAndServe()
}
