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
	gt.SetServiceDesc(&desc)

	eid := uint16(snowflake.GenIDU32())
	gt.entityID = entity.FixedEntityID(0, entity.InstID(eid), entity.ID(eid))

	gt.EntityMgr.RegisterEntity(nil, gt.entityID, gt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go gt.Run(&wg)
	wg.Wait()

	gt.init()

	return gt
}

func (gt *Gate) init() {
	ex.ConnProxy(nets.WithPacketHandler(gt))

	gt.EntityMgr.RemoteEntityManager = ex.ProxyClient

	client := proxypb.NewProxySvcServiceClient()

	rsp, err := client.RegisterInst(gt.Context(), &proxypb.RegisterInstRequest{
		InstID: uint32(gt.EntityID().InstID()),
	})
	if err != nil {
		xlog.Errorln("error:", err)
		panic(err)
	}

	// 更新entityID
	proxyID := entity.ProxyID(rsp.GetProxyID())
	gt.entityID = entity.FixedEntityID(proxyID, gt.Service.EntityID().InstID(), gt.Service.EntityID().ID())
}

type Gate struct {
	entity.Service

	cltMtx sync.Mutex
	cltMap clientMap

	entityID entity.EntityID

	*bs.Server
}

type clientMap map[entity.ID]*codec.PacketReadWriter

// // RegisterEntity register serive
func (gt *Gate) RegisterEntity(eid entity.EntityID, prw *codec.PacketReadWriter) {
	gt.cltMtx.Lock()
	defer gt.cltMtx.Unlock()

	gt.cltMap[eid.ID()] = prw
}

// GateService
type GateService struct {
	entity.Service
}

// RegisterClient
func (gt *Gate) RegisterClient(c entity.Context, req *gatepb.RegisterClientRequest) (*gatepb.RegisterClientResponse, error) {

	id := gt.NewEntityID()

	gt.RegisterEntity(id, c.Packet().Src)

	// client := proxypb.NewProxySvcServiceClient()
	// rsp, err := client.RegisterEntity(c, &proxypb.RegisterEntityRequest{EntityID: req.EntityID})
	// if err != nil {
	// 	return nil, err
	// }
	// xlog.Debugln("register proxy entity resp", rsp.String())
	return &gatepb.RegisterClientResponse{
		EntityID: uint64(id),
	}, nil
}

// UnregisterClient
func (gt *Gate) UnregisterClient(c entity.Context, req *gatepb.RegisterClientRequest) error {
	return nil
}

// Ping
func (gt *Gate) Ping(c entity.Context, req *gatepb.PingPong) (*gatepb.PingPong, error) {
	return nil, nil
}

// func (gt *Gate) RegisterEntityToProxy(eid entity.EntityID) error {
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

func (gt *Gate) NewEntityID() entity.EntityID {
	return entity.NewEntityID(gt.EntityID().ProxyID(), gt.EntityID().InstID())
}

func (gt *Gate) EntityID() entity.EntityID {
	return gt.entityID
}

func (gt *Gate) Serve() error {
	return gt.Server.ListenAndServe()
}
