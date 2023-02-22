package main

import (
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

type clientMap map[string]*codec.PacketReadWriter

var cltMap clientMap = make(clientMap)

// // RegisterEntity register serive
func RegisterEntity(sid *bbq.EntityID, prw *codec.PacketReadWriter) {
	cltMap[sid.ID] = prw
}

// GateService
type GateService struct {
	entity.Service
}

// RegisterClient
func (gs *GateService) RegisterClient(c entity.Context, req *gatepb.RegisterClientRequest) (*gatepb.RegisterClientResponse, error) {

	req.EntityID.ProxyID = Inst.ProxyID

	RegisterEntity(req.EntityID, c.Packet().Src)

	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())
	rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: req.EntityID})
	if err != nil {
		return nil, err
	}
	xlog.Println("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{EntityID: req.EntityID}, nil
}

// UnregisterClient
func (gs *GateService) UnregisterClient(c entity.Context, req *gatepb.RegisterClientRequest) {

}

// Ping
func (gs *GateService) Ping(c entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

type Gate struct {
	entity.Entity
	ProxyID string
}

func NewGate() *Gate {

	conf.Init("gate.yaml")

	entity.ProxyRegister = &RegisterProxy{}
	entity.NewEntityID = &EntityIDGenerator{}

	gm := &Gate{}
	eid := snowflake.GenUUID()

	entity.RegisterEntity(nil, &bbq.EntityID{ID: eid}, gm)

	go gm.Run()

	return gm
}

var Inst = NewGate()

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid *bbq.EntityID) error {
	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterEntity(Inst.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
	if err != nil {
		return err
	}

	xlog.Println("register proxy entity resp")

	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterService(Inst.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Println("register proxy service resp")

	return nil
}

type EntityIDGenerator struct {
}

func (n *EntityIDGenerator) NewEntityID(tn string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: tn, ProxyID: Inst.ProxyID}
}
