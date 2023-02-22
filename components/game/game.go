package game

import (
	"time"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

func Init() {

	ex.ConnProxy(nets.WithPacketHandler(NewGamePacketHandler()))

	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())

	rsp, err := client.RegisterInst(Inst.Context(), &proxypb.RegisterInstRequest{
		InstID: Inst.EntityID().ID,
	})
	if err != nil {
		panic(err)
	}

	Inst.ProxyID = rsp.ProxyID
}

func Run() {
	for {
		xlog.Info("Run Game")
		sleepTime := 50
		for {
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}
}

type Game struct {
	entity.Entity

	ProxyID string
}

func NewGame() *Game {

	conf.Init("game.yaml")

	entity.ProxyRegister = &RegisterProxy{}
	entity.NewEntityID = &EntityIDGenerator{}

	gm := &Game{}
	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, &entity.EntityID{ID: eid}, gm)

	go gm.Run()

	return gm
}

var Inst = NewGame()

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid *entity.EntityID) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterEntity(Inst.Context(), &proxypb.RegisterEntityRequest{EntityID: entity.ToPBEntityID(eid)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy entity resp")
	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterService(Inst.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy service resp")

	return nil
}

type EntityIDGenerator struct {
}

func (n *EntityIDGenerator) NewEntityID(tn entity.TypeName) *entity.EntityID {
	return &entity.EntityID{ID: snowflake.GenUUID(), Type: tn, ProxyID: Inst.ProxyID}
}
