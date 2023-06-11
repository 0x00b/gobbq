package client

import (
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/snowflake"
)

type Client struct {
	entity.IEntity

	Gate *nets.Client

	EntityMgr *entity.EntityManager

	entityID entity.EntityID
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

	client.EntityMgr.RemoteEntityManager = gate

	client.EntityMgr.RegisterEntityDesc(sd, ss)

	eid := snowflake.GenIDU32()
	client.entityID = entity.FixedEntityID(0, 0, entity.ID(eid))

	client.IEntity, err = client.EntityMgr.NewEntity(nil, client.entityID, sd.TypeName)
	if err != nil {
		panic(err)
	}

	client.Gate = gate

	gateSvc := gatepb.NewGateServiceClient()
	// go func() {
	// 	client.Run()

	// 	// unregister
	// 	gateSvc.UnregisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: client.EntityID()})
	// }()

	// time.Sleep(1 * time.Second)

	rsp, err := gateSvc.RegisterClient(client.Context(), &gatepb.RegisterClientRequest{
		// EntityID: uint64(client.IEntity.EntityID().ID()),
	})

	if err != nil {
		panic(err)
	}

	newID := entity.EntityID(rsp.GetEntityID())
	client.EntityMgr.ReplaceEntityID(client.Context(), client.entityID, newID)
	client.entityID = newID

	return client
}

// EntityID 重写EntityID
func (c *Client) EntityID() entity.EntityID {
	return c.entityID
}
