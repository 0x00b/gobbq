package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
	capi "github.com/hashicorp/consul/api"
)

type Game struct {
	entity.Entity

	// 可以重写这个接口，自己规划EntityID中的ID部分，比如用来承载UID等
	IDGenerator

	EntityMgr *entity.EntityManager

	// todo
	DB db.IDatabase

	consul *capi.Client
}

type IDGenerator interface {
	GenID() entity.ID
}

func NewGame(opts ...Option) *Game {

	InitConfig()

	gm := &Game{
		EntityMgr:   entity.NewEntityManager(),
		IDGenerator: &defaultIDGener{},
	}
	gm.EntityMgr.ProxyRegister = gm
	gm.EntityMgr.EntityIDGenerator = gm

	for _, opt := range opts {
		opt(gm)
	}

	eid := gm.IDGenerator.GenID()

	desc := entity.EntityDesc{}
	desc.EntityImpl = gm
	desc.EntityMgr = gm.EntityMgr
	entity.SetEntityDesc(gm, &desc)

	temp := entity.FixedEntityID(0, 0, eid)

	gm.EntityMgr.RegisterEntity(nil, temp, gm)

	entity.Run(gm)

	gm.init(temp)

	return gm
}

func (g *Game) init(old entity.EntityID) {

	cfg := CFG.Proxy
	if !CFG.LocalProxy {
		var err error
		g.consul, err = capi.NewClient(capi.DefaultConfig())
		if err != nil {
			panic(err)
		}

		_, svcs, err := g.consul.Agent().AgentHealthServiceByName("gobbqproxy")
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
	ex.ConnProxy(cfg, nets.WithPacketHandler(g.EntityMgr), nets.WithConnCallback(g))

	g.EntityMgr.Proxy = ex.ProxyClient

	client := proxypb.NewProxySvcClient()

	rsp, err := client.RegisterInst(g.Context(), &proxypb.RegisterInstRequest{})
	if err != nil {
		panic(err)
	}

	proxyID := entity.ProxyID(rsp.GetProxyID())
	instID := entity.InstID(rsp.GetInstID())
	new := entity.FixedEntityID(proxyID, instID, g.Entity.EntityID().ID())
	g.EntityMgr.ReplaceEntityID(old, new)

	g.AddTimer(10*time.Second, func() bool {
		client.Ping(g.Context(), &proxypb.PingPong{})
		return true
	})
}

func (g *Game) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxyEtyClient(entity.FixedEntityID(g.EntityID().ProxyID(), 0, 0))

	_, err := client.RegisterService(g.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debug("register proxy service succ:", svcName)

	return nil
}

// NewEntityID 如果没有特殊规划，可以使用这个生成entity id
func (g *Game) NewEntityID() entity.EntityID {
	return entity.FixedEntityID(g.EntityID().ProxyID(), g.EntityID().InstID(), g.IDGenerator.GenID())
}

func (g *Game) Serve() {
	for {
		xlog.Debug("Run Game")
		sleepTime := 50
		for {
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}
}
