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

func Init() {

	ex.ConnProxy(nets.WithPacketHandler(NewGamePacketHandler()))

	client := proxypb.NewProxySvcServiceClient(ex.ProxyClient.GetPacketReadWriter())

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

	entity.RegisterEntity(nil, &bbq.EntityID{ID: eid, ProxyID: eid}, gm)

	go gm.Run()

	return gm
}

var Inst = NewGame()

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid *bbq.EntityID) error {

	client := proxypb.NewProxySvcServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterEntity(Inst.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
	if err != nil {
		return err
	}

	xlog.Println("register proxy entity resp")
	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxySvcServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterService(Inst.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy service resp")

	return nil
}

type EntityIDGenerator struct {
}

func (n *EntityIDGenerator) NewEntityID(typeName string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: typeName, ProxyID: Inst.ProxyID}
}
