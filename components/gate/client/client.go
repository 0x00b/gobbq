package client

import (
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/snowflake"
)

type Client struct {
	entity.Service

	Gate *nets.Client
}

func NewClient() *Client {

	cfg := conf.C.Gate.Inst[0]
	gate, err := nets.Connect(
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(&ClientPacketHandler{}))
	if err != nil {
		panic(err)
	}

	client := &Client{}

	eid := &entity.EntityID{ID: snowflake.GenUUID()}
	entity.RegisterEntity(nil, eid, client)

	client.Gate = gate

	gateSvc := gatepb.NewGateServiceClient(gate.GetPacketReadWriter())
	go func() {
		client.Run()

		// unregister
		gateSvc.UnregisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: entity.ToPBEntityID(client.EntityID())})
	}()

	rsp, err := gateSvc.RegisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: entity.ToPBEntityID(client.EntityID())})

	if err != nil {
		panic(err)
	}

	eid.ProxyID = rsp.EntityID.ProxyID

	return client
}

type ClientPacketHandler struct {
	entity.MethodPacketHandler
}
