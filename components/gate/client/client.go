package client

import (
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
)

type Client struct {
	entity.IEntity

	Gate *nets.Client
}

func NewClient(sd *entity.EntityDesc, ss entity.IEntity, intercepter ...entity.ServerInterceptor) *Client {

	cfg := conf.C.Gate.Inst[0]
	gate, err := nets.Connect(
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(&ClientPacketHandler{}))
	if err != nil {
		panic(err)
	}

	client := &Client{}

	entity.Manager.RegisterEntity(sd, ss)

	eid := &bbq.EntityID{ID: snowflake.GenUUID(), Type: sd.TypeName}

	client.IEntity, err = entity.NewEntity(nil, eid)
	if err != nil {
		panic(err)
	}

	client.Gate = gate

	gateSvc := gatepb.NewGateServiceClient(gate.GetPacketReadWriter())
	go func() {
		client.Run()

		// unregister
		gateSvc.UnregisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: client.EntityID()})
	}()

	rsp, err := gateSvc.RegisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: client.EntityID()})

	if err != nil {
		panic(err)
	}

	client.EntityID().ProxyID = rsp.EntityID.ProxyID

	return client
}

type ClientPacketHandler struct {
	entity.MethodPacketHandler
}
