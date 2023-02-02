package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
)

type clientMap map[entity.EntityID]*codec.PacketReadWriter

var cltMap clientMap = make(clientMap)

// // RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	cltMap[sid] = prw
}

// GateService
type GateService struct {
	entity.Service
}

// RegisterClient
func (gs *GateService) RegisterClient(c *entity.Context, req *gatepb.RegisterClientRequest, ret func(*gatepb.RegisterClientResponse, error)) {
	RegisterEntity(entity.EntityID(req.EntityID), c.Packet().Src)
}

// UnregisterClient
func (gs *GateService) UnregisterClient(c *entity.Context, req *gatepb.RegisterClientRequest) {

}

// Ping
func (gs *GateService) Ping(c *entity.Context, req *gatepb.PingPong, ret func(*gatepb.PingPong, error)) {

}

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid entity.EntityID) error {
	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterEntity(nil, &proxypb.RegisterEntityRequest{EntityID: string(eid)},
		func(c *entity.Context, rsp *proxypb.RegisterEntityResponse) {
			fmt.Println("register proxy entity resp")
		},
	)

	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterService(nil, &proxypb.RegisterServiceRequest{ServiceName: string(svcName)},
		func(c *entity.Context, rsp *proxypb.RegisterServiceResponse) {

			fmt.Println("register proxy service resp")
		},
	)

	return nil
}
