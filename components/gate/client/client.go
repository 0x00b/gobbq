package client

import (
	"time"

	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

type Client struct {
	entity.IEntity

	Gate *nets.Client

	EntityMgr *entity.EntityManager
}

func NewClient(sd *entity.EntityDesc, ss entity.IEntity, intercepter ...entity.ServerInterceptor) *Client {

	client := &Client{
		EntityMgr: entity.NewEntityManager(),
	}

	cfg := conf.C.Gate.Inst[0]
	gate, err := nets.Connect(
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(client.EntityMgr))
	if err != nil {
		panic(err)
	}

	client.EntityMgr.Proxy = gate

	client.EntityMgr.RegisterEntityDesc(sd, ss)

	// 临时的
	eid := entity.FixedEntityID(0, 0, entity.ID(entity.GenIDU32()))

	client.IEntity, err = client.EntityMgr.NewEntity(nil, eid, sd.TypeName)
	if err != nil {
		panic(err)
	}

	client.Gate = gate

	gateSvc := gatepb.NewGateServiceClient()
	// secure.GO( func() {
	// 	client.Run()

	// 	// unregister
	// 	gateSvc.UnregisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: client.EntityID()})
	// })

	// time.Sleep(1 * time.Second)

	rsp, err := gateSvc.RegisterClient(client.Context(), &gatepb.RegisterClientRequest{
		// EntityID: uint64(client.IEntity.EntityID().ID()),
	})

	if err != nil {
		panic(err)
	}

	newID := entity.EntityID(rsp.GetEntityID())
	client.EntityMgr.ReplaceEntityID(eid, newID)

	client.AddTimer(10*time.Second, func() {
		gateSvc.Ping(client.Context(), &gatepb.PingPong{})
	})

	return client
}
