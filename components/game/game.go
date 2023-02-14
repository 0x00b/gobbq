package game

import (
	"time"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

func Init() {

	ex.ConnProxy(nets.WithPacketHandler(NewGamePacketHandler()))
	entity.ProxyRegister = &RegisterProxy{}
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
}

func NewGame() *Game {
	gm := &Game{}
	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, entity.EntityID(eid), gm)

	go gm.Run()

	return gm
}

var Inst = NewGame()

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid entity.EntityID) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	_, err := client.RegisterEntity(Inst.Context(), &proxypb.RegisterEntityRequest{EntityID: string(eid)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy entity resp")
	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	_, err := client.RegisterService(Inst.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy service resp")

	return nil
}
