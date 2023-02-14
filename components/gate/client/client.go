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

	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, entity.EntityID(eid), client)

	client.Gate = gate

	gateSvc := gatepb.NewGateServiceClient(gate)
	go func() {
		client.Run()

		// unregister
		gateSvc.UnregisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: eid})
	}()

	_, err = gateSvc.RegisterClient(client.Context(), &gatepb.RegisterClientRequest{EntityID: eid})

	if err != nil {
		panic(err)
	}

	return client
}

type ClientPacketHandler struct {
	entity.MethodPacketHandler
}
