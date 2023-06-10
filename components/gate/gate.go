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

	gt.EntityMgr.ProxyRegister = gt
	gt.EntityMgr.EntityIDGenerator = gt

	desc := gatepb.GateServiceDesc
	desc.EntityImpl = gt
	desc.EntityMgr = gt.EntityMgr
	gt.SetDesc(&desc)

	eid := &bbq.EntityID{ID: snowflake.GenUUID(), Type: gatepb.GateServiceDesc.TypeName}

	gt.EntityMgr.RegisterEntity(nil, eid, gt)

	go gt.Run()

	gt.init()

	return gt
}

func (gt *Gate) init() {
	ex.ConnProxy(nets.WithPacketHandler(gt))

	gt.EntityMgr.RemoteEntityManager = ex.ProxyClient

	client := proxypb.NewProxySvcServiceClient()

	rsp, err := client.RegisterInst(gt.Context(), &proxypb.RegisterInstRequest{
		InstID: gt.EntityID().ID,
	})
	if err != nil {
		xlog.Errorln("error:", err)
		panic(err)
	}

	gt.EntityID().ProxyID = rsp.ProxyID
}

type Gate struct {
	entity.Service

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

	gt.RegisterEntity(req.EntityID, c.Packet().Src)

	// client := proxypb.NewProxySvcServiceClient()
	// rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: req.EntityID})
	// if err != nil {
	// 	return nil, err
	// }
	// xlog.Debugln("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{EntityID: gt.EntityID()}, nil
}

// UnregisterClient
func (gt *Gate) UnregisterClient(c entity.Context, req *gatepb.RegisterClientRequest) error {
	return nil
}

// Ping
func (gt *Gate) Ping(c entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

// func (gt *Gate) RegisterEntityToProxy(eid *bbq.EntityID) error {
// 	client := proxypb.NewProxySvcServiceClient()

// 	_, err := client.RegisterEntity(gt.Context(), &proxypb.RegisterEntityRequest{EntityID: eid})
// 	if err != nil {
// 		return err
// 	}

// 	xlog.Debugln("register proxy entity resp")

// 	return nil
// }

func (gt *Gate) RegisterServiceToProxy(svcName string) error {

	client := proxypb.NewProxySvcServiceClient()

	_, err := client.RegisterService(gt.Context(), &proxypb.RegisterServiceRequest{ServiceName: string(svcName)})
	if err != nil {
		return err
	}

	xlog.Debugln("register proxy service resp")

	return nil
}

func (gt *Gate) NewEntityID(typeName string) *bbq.EntityID {
	return &bbq.EntityID{ID: snowflake.GenUUID(), Type: typeName, ProxyID: gt.EntityID().ProxyID, InstID: gt.EntityID().ID}
}

func (gt *Gate) Serve() error {
	return gt.Server.ListenAndServe()
}
