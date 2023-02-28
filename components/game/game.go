package game

import (
	"time"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

type Game struct {
	entity.Entity

	ProxyID string

	EntityMgr *entity.EntityManager
}

func NewGame() *Game {

	conf.Init("game.yaml")

	gm := &Game{
		EntityMgr: entity.NewEntityManager(),
	}
	gm.EntityMgr.ProxyRegister = gm
	gm.EntityMgr.EntityIDGenerator = gm

	eid := snowflake.GenUUID()

	desc := entity.EntityDesc{}
	desc.EntityImpl = gm
	desc.EntityMgr = gm.EntityMgr
	gm.SetDesc(&desc)

	gm.EntityMgr.RegisterEntity(nil, &bbq.EntityID{ID: eid, ProxyID: eid}, gm)

	go gm.Run()

	gm.init()

	return gm
}

func (g *Game) init() {

	ex.ConnProxy(nets.WithPacketHandler(g.EntityMgr))

	g.EntityMgr.RemoteEntityManager = ex.ProxyClient

	client := proxypb.NewProxySvcServiceClient()

	rsp, err := client.RegisterInst(g.Context(), &proxypb.RegisterInstRequest{
		InstID: g.EntityID().ID,
	})
	if err != nil {
		panic(err)
	}

	g.ProxyID = rsp.ProxyID
}

func (g *Game) RegisterEntityToProxy(eid *bbq.EntityID) error {

	client := proxypb.NewProxySvcServiceClient()

	_, err := client.RegisterEntity(g.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
	if err != nil {
		return err
	}

	xlog.Debug("register proxy entity resp")
	return nil
}

func (g *Game) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxySvcServiceClient()

	_, err := client.RegisterService(g.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debug("register proxy service resp")

	return nil
}

func (g *Game) NewEntityID(typeName string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: typeName, ProxyID: g.ProxyID}
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
