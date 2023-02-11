package main

import (
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

type clientMap map[entity.EntityID]*codec.PacketReadWriter

var cltMap clientMap = make(clientMap)

// // RegisterEntity register serive
func RegisterEntity(sid entity.EntityID, prw *codec.PacketReadWriter) {
	cltMap[sid] = prw
}

// GateService
type GateService struct {
	entity.Entity
}

// RegisterClient
func (gs *GateService) RegisterClient(c *entity.Context, req *gatepb.RegisterClientRequest, ret func(*gatepb.RegisterClientResponse, error)) {
	eid := snowflake.GenUUID()

	RegisterEntity(entity.EntityID(eid), c.Packet().Src)

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)
	client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: string(eid)},
		func(c *entity.Context, rsp *proxypb.RegisterEntityResponse) {
			xlog.Println("register proxy entity resp")
		},
	)
	ret(&gatepb.RegisterClientResponse{EntityID: eid}, nil)
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
			xlog.Println("register proxy entity resp")
		},
	)

	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterService(nil, &proxypb.RegisterServiceRequest{ServiceName: string(svcName)},
		func(c *entity.Context, rsp *proxypb.RegisterServiceResponse) {

			xlog.Println("register proxy service resp")
		},
	)

	return nil
}
