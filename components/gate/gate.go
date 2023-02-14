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
	entity.Service
}

// RegisterClient
func (gs *GateService) RegisterClient(c *entity.Context, req *gatepb.RegisterClientRequest) (*gatepb.RegisterClientResponse, error) {

	RegisterEntity(entity.EntityID(req.EntityID), c.Packet().Src)

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)
	rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: string(req.EntityID)})
	if err != nil {
		return nil, err
	}
	xlog.Println("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{EntityID: req.EntityID}, nil
}

// UnregisterClient
func (gs *GateService) UnregisterClient(c *entity.Context, req *gatepb.RegisterClientRequest) {

}

// Ping
func (gs *GateService) Ping(c *entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

type Gate struct {
	entity.Entity
}

func NewGate() *Gate {
	gm := &Gate{}
	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, entity.EntityID(eid), gm)

	go gm.Run()

	return gm
}

var Inst = NewGate()

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
