package main

import (
	"sync"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

func NewGate() *Gate {

	conf.Init("gate.yaml")

	gt := &Gate{
		Server: bs.NewServer(),
		cltMap: make(clientMap),
	}

	gt.Server.EntityMgr.ProxyRegister = gt
	gt.Server.EntityMgr.EntityIDGenerator = gt

	eid := snowflake.GenUUID()

	gt.Server.EntityMgr.RegisterEntity(nil, &bbq.EntityID{ID: eid}, gt)

	go gt.Run()

	gt.init()

	return gt
}

type Gate struct {
	entity.Service
	ProxyID string

	cltMtx sync.Mutex
	cltMap clientMap

	*bs.Server
}

type clientMap map[string]*codec.PacketReadWriter

// // RegisterEntity register serive
func (gt *Gate) RegisterEntity(eid *bbq.EntityID, prw *codec.PacketReadWriter) {
	gt.cltMtx.Lock()
	defer gt.cltMtx.Unlock()

	gt.cltMap[eid.ID] = prw
}

// GateService
type GateService struct {
	entity.Service
}

// RegisterClient
func (gt *Gate) RegisterClient(c entity.Context, req *gatepb.RegisterClientRequest) (*gatepb.RegisterClientResponse, error) {

	req.EntityID.ProxyID = gt.ProxyID

	gt.RegisterEntity(req.EntityID, c.Packet().Src)

	client := proxypb.NewProxySvcServiceClient(gt.EntityMgr, ex.ProxyClient.GetPacketReadWriter())
	rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: req.EntityID})
	if err != nil {
		return nil, err
	}
	xlog.Debugln("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{EntityID: req.EntityID}, nil
}

// UnregisterClient
func (gt *Gate) UnregisterClient(c entity.Context, req *gatepb.RegisterClientRequest) error {
	return nil
}

// Ping
func (gt *Gate) Ping(c entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

func (gt *Gate) RegisterEntityToProxy(eid *bbq.EntityID) error {
	client := proxypb.NewProxySvcServiceClient(gt.EntityMgr, ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterEntity(gt.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
	if err != nil {
		return err
	}

	xlog.Debugln("register proxy entity resp")

	return nil
}

func (gt *Gate) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxySvcServiceClient(gt.EntityMgr, ex.ProxyClient.GetPacketReadWriter())

	_, err := client.RegisterService(gt.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debugln("register proxy service resp")

	return nil
}

func (gt *Gate) NewEntityID(typeName string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: typeName, ProxyID: gt.ProxyID}
}

func (gt *Gate) init() {
	ex.ConnProxy(nets.WithPacketHandler(gt))
	client := proxypb.NewProxySvcServiceClient(gt.EntityMgr, ex.ProxyClient.GetPacketReadWriter())

	rsp, err := client.RegisterInst(gt.Context(), &proxypb.RegisterInstRequest{
		InstID: gt.EntityID().ID,
	})
	if err != nil {
		panic(err)
	}

	gt.ProxyID = rsp.ProxyID
}

func (gt *Gate) Serve() error {
	return gt.Server.ListenAndServe()
}
