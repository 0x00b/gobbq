package game

import (
	"sync"
	"time"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

type Game struct {
	entity.Entity

	EntityMgr *entity.EntityManager

	entityID entity.EntityID
}

func NewGame() *Game {

	conf.Init("game.yaml")

	gm := &Game{
		EntityMgr: entity.NewEntityManager(),
	}
	gm.EntityMgr.ProxyRegister = gm
	gm.EntityMgr.EntityIDGenerator = gm

	eid := uint16(snowflake.GenIDU32())

	desc := entity.EntityDesc{}
	desc.EntityImpl = gm
	desc.EntityMgr = gm.EntityMgr
	gm.SetEntityDesc(&desc)

	gm.entityID = entity.FixedEntityID(0, entity.InstID(eid), entity.ID(eid))

	gm.EntityMgr.RegisterEntity(nil, gm.entityID, gm)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go gm.Run(&wg)
	wg.Wait()

	gm.init()

	return gm
}

// EntityID 重写
func (g *Game) EntityID() entity.EntityID {
	return g.entityID
}

func (g *Game) init() {

	ex.ConnProxy(nets.WithPacketHandler(g.EntityMgr))

	g.EntityMgr.RemoteEntityManager = ex.ProxyClient

	client := proxypb.NewProxySvcServiceClient()

	rsp, err := client.RegisterInst(g.Context(), &proxypb.RegisterInstRequest{
		InstID: uint32(g.EntityID().InstID()),
	})
	if err != nil {
		panic(err)
	}

	proxyID := entity.ProxyID(rsp.GetProxyID())
	g.entityID = entity.FixedEntityID(proxyID, g.Entity.EntityID().InstID(), g.Entity.EntityID().ID())
}

// func (g *Game) RegisterEntityToProxy(eid entity.EntityID) error {

// 	client := proxypb.NewProxySvcServiceClient()

// 	_, err := client.RegisterEntity(g.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
// 	if err != nil {
// 		return err
// 	}

// 	xlog.Debug("register proxy entity resp")
// 	return nil
// }

func (g *Game) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxySvcServiceClient()

	_, err := client.RegisterService(g.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debug("register proxy service resp")

	return nil
}

func (g *Game) NewEntityID() entity.EntityID {
	return entity.NewEntityID(g.EntityID().ProxyID(), g.EntityID().InstID())
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
