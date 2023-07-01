// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package proxypb

import (
	"errors"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
	// proxypb "github.com/0x00b/gobbq/components/proxy/proxypb"
)

var _ = snowflake.GenUUID()

func RegisterProxyEtyEntity(etyMgr *entity.EntityManager, impl ProxyEtyEntity) {
	etyMgr.RegisterEntityDesc(&ProxyEtyEntityDesc, impl)
}

func NewProxyEtyClient(eid entity.EntityID) *ProxyEty {
	t := &ProxyEty{
		EntityID: eid,
	}
	return t
}

func NewProxyEty(c entity.Context) *ProxyEty {
	etyMgr := entity.GetEntityMgr(c)
	return NewProxyEtyWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewProxyEtyWithID(c entity.Context, id entity.EntityID) *ProxyEty {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, ProxyEtyEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &ProxyEty{
		EntityID: id,
	}

	return t
}

type ProxyEty struct {
	EntityID entity.EntityID
}

func (t *ProxyEty) RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ProxyEtyEntityDesc.TypeName
	pkt.Header.Method = "RegisterProxy"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(RegisterProxyResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*RegisterProxyResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *ProxyEty) SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ProxyEtyEntityDesc.TypeName
	pkt.Header.Method = "SyncService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(SyncServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*SyncServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *ProxyEty) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ProxyEtyEntityDesc.TypeName
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ProxyEtyEntity
type ProxyEtyEntity interface {
	entity.IEntity

	// RegisterProxy
	RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error)

	// SyncService
	SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error)

	// Ping
	Ping(c entity.Context, req *PingPong) (*PingPong, error)
}

func _ProxyEtyEntity_RegisterProxy_Handler(svc any, ctx entity.Context, in *RegisterProxyRequest, interceptor entity.ServerInterceptor) (*RegisterProxyResponse, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).RegisterProxy(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/RegisterProxy",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).RegisterProxy(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterProxyResponse), err

}

func _ProxyEtyEntity_RegisterProxy_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_RegisterProxy_Handler(svc, ctx, in.(*RegisterProxyRequest), interceptor)

}

func _ProxyEtyEntity_RegisterProxy_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterProxyRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_RegisterProxy_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _ProxyEtyEntity_SyncService_Handler(svc any, ctx entity.Context, in *SyncServiceRequest, interceptor entity.ServerInterceptor) (*SyncServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).SyncService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/SyncService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).SyncService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SyncServiceResponse), err

}

func _ProxyEtyEntity_SyncService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_SyncService_Handler(svc, ctx, in.(*SyncServiceRequest), interceptor)

}

func _ProxyEtyEntity_SyncService_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SyncServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_SyncService_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _ProxyEtyEntity_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).Ping(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*PingPong), err

}

func _ProxyEtyEntity_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_Ping_Handler(svc, ctx, in.(*PingPong), interceptor)

}

func _ProxyEtyEntity_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_Ping_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

var ProxyEtyEntityDesc = entity.EntityDesc{
	TypeName:    "proxypb.ProxyEtyEntity",
	HandlerType: (*ProxyEtyEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterProxy": {
			MethodName:   "RegisterProxy",
			Handler:      _ProxyEtyEntity_RegisterProxy_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_RegisterProxy_Local_Handler,
		},

		"SyncService": {
			MethodName:   "SyncService",
			Handler:      _ProxyEtyEntity_SyncService_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_SyncService_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxyEtyEntity_Ping_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}

func RegisterProxySvcService(etyMgr *entity.EntityManager, impl ProxySvcService) {
	etyMgr.RegisterService(&ProxySvcServiceDesc, impl)
}

func NewProxySvcClient() *ProxySvc {
	t := &ProxySvc{}
	return t
}

type ProxySvc struct {
}

func (t *ProxySvc) RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = ProxySvcServiceDesc.TypeName
	pkt.Header.Method = "RegisterInst"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(RegisterInstResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*RegisterInstResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *ProxySvc) RegisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = ProxySvcServiceDesc.TypeName
	pkt.Header.Method = "RegisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *ProxySvc) UnregisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = ProxySvcServiceDesc.TypeName
	pkt.Header.Method = "UnregisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *ProxySvc) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = ProxySvcServiceDesc.TypeName
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ProxySvcService
type ProxySvcService interface {
	entity.IService

	// RegisterInst
	RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error)

	// RegisterService
	RegisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error)

	// UnregisterService
	UnregisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error)

	// Ping
	Ping(c entity.Context, req *PingPong) (*PingPong, error)
}

func _ProxySvcService_RegisterInst_Handler(svc any, ctx entity.Context, in *RegisterInstRequest, interceptor entity.ServerInterceptor) (*RegisterInstResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).RegisterInst(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/RegisterInst",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).RegisterInst(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterInstResponse), err

}

func _ProxySvcService_RegisterInst_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_RegisterInst_Handler(svc, ctx, in.(*RegisterInstRequest), interceptor)

}

func _ProxySvcService_RegisterInst_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterInstRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_RegisterInst_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _ProxySvcService_RegisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).RegisterService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/RegisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).RegisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterServiceResponse), err

}

func _ProxySvcService_RegisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_RegisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), interceptor)

}

func _ProxySvcService_RegisterService_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_RegisterService_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _ProxySvcService_UnregisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).UnregisterService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/UnregisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).UnregisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterServiceResponse), err

}

func _ProxySvcService_UnregisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_UnregisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), interceptor)

}

func _ProxySvcService_UnregisterService_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_UnregisterService_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _ProxySvcService_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).Ping(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*PingPong), err

}

func _ProxySvcService_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_Ping_Handler(svc, ctx, in.(*PingPong), interceptor)

}

func _ProxySvcService_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_Ping_Handler(svc, ctx, in, interceptor)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	if err != nil {
		npkt.Header.ErrCode = 1
		npkt.Header.ErrMsg = err.Error()

		npkt.WriteBody(nil)
	} else {
		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			xlog.Errorln("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

var ProxySvcServiceDesc = entity.EntityDesc{
	TypeName:    "proxypb.ProxySvcService",
	HandlerType: (*ProxySvcService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterInst": {
			MethodName:   "RegisterInst",
			Handler:      _ProxySvcService_RegisterInst_Remote_Handler,
			LocalHandler: _ProxySvcService_RegisterInst_Local_Handler,
		},

		"RegisterService": {
			MethodName:   "RegisterService",
			Handler:      _ProxySvcService_RegisterService_Remote_Handler,
			LocalHandler: _ProxySvcService_RegisterService_Local_Handler,
		},

		"UnregisterService": {
			MethodName:   "UnregisterService",
			Handler:      _ProxySvcService_UnregisterService_Remote_Handler,
			LocalHandler: _ProxySvcService_UnregisterService_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxySvcService_Ping_Remote_Handler,
			LocalHandler: _ProxySvcService_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}
