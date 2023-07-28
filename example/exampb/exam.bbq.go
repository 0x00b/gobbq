// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package exampb

import (
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"

	// exampb "github.com/0x00b/gobbq/example/exampb"

)

var _ = snowflake.GenUUID()

func RegisterEchoService(etyMgr *entity.EntityManager, impl EchoService) {
	etyMgr.RegisterService(&EchoServiceDesc, impl)
}

func NewEchoClient() *Echo {
	t := &Echo{}
	return t
}

type Echo struct {
}

func (t *Echo) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = EchoServiceDesc.TypeName
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// 如果是LocalCall，由local内部关闭chan
	isLocalCall := false
	var chanRsp chan any = make(chan any)
	defer func() {
		if !isLocalCall {
			close(chanRsp)
		}
	}()

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)

	isLocalCall = err == nil

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
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(SayHelloResponse)
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
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// EchoService
type EchoService interface {
	entity.IService

	// SayHello
	SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _EchoService_SayHello_Handler(svc any, ctx entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {
		return svc.(EchoService).SayHello(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoService/SayHello",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(EchoService).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SayHelloResponse), err

}

func _EchoService_SayHello_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _EchoService_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), interceptor)

}

func _EchoService_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _EchoService_SayHello_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

var EchoServiceDesc = entity.EntityDesc{
	TypeName:    "exampb.EchoService",
	HandlerType: (*EchoService)(nil),
	Methods: map[string]entity.MethodDesc{

		"SayHello": {
			MethodName:   "SayHello",
			Handler:      _EchoService_SayHello_Remote_Handler,
			LocalHandler: _EchoService_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}

func RegisterEchoEtyEntity(etyMgr *entity.EntityManager, impl EchoEtyEntity) {
	etyMgr.RegisterEntityDesc(&EchoEtyEntityDesc, impl)
}

func NewEchoEtyClient(eid entity.EntityID) *EchoEty {
	t := &EchoEty{
		EntityID: eid,
	}
	return t
}

func NewEchoEty(c entity.Context) (*EchoEty, error) {
	etyMgr := entity.GetEntityMgr(c)
	return NewEchoEtyWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewEchoEtyWithID(c entity.Context, id entity.EntityID) (*EchoEty, error) {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, EchoEtyEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil, err
	}
	t := &EchoEty{
		EntityID: id,
	}

	return t, nil
}

type EchoEty struct {
	EntityID entity.EntityID
}

func (t *EchoEty) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// 如果是LocalCall，由local内部关闭chan
	isLocalCall := false
	var chanRsp chan any = make(chan any)
	defer func() {
		if !isLocalCall {
			close(chanRsp)
		}
	}()

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)

	isLocalCall = err == nil

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
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(SayHelloResponse)
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
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// EchoEtyEntity
type EchoEtyEntity interface {
	entity.IEntity

	// SayHello
	SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _EchoEtyEntity_SayHello_Handler(svc any, ctx entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {
		return svc.(EchoEtyEntity).SayHello(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoEtyEntity/SayHello",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(EchoEtyEntity).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SayHelloResponse), err

}

func _EchoEtyEntity_SayHello_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _EchoEtyEntity_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), interceptor)

}

func _EchoEtyEntity_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _EchoEtyEntity_SayHello_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

var EchoEtyEntityDesc = entity.EntityDesc{
	TypeName:    "exampb.EchoEtyEntity",
	HandlerType: (*EchoEtyEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"SayHello": {
			MethodName:   "SayHello",
			Handler:      _EchoEtyEntity_SayHello_Remote_Handler,
			LocalHandler: _EchoEtyEntity_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}

func RegisterEchoSvc2Service(etyMgr *entity.EntityManager, impl EchoSvc2Service) {
	etyMgr.RegisterService(&EchoSvc2ServiceDesc, impl)
}

func NewEchoSvc2Client() *EchoSvc2 {
	t := &EchoSvc2{}
	return t
}

type EchoSvc2 struct {
}

func (t *EchoSvc2) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = EchoSvc2ServiceDesc.TypeName
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// 如果是LocalCall，由local内部关闭chan
	isLocalCall := false
	var chanRsp chan any = make(chan any)
	defer func() {
		if !isLocalCall {
			close(chanRsp)
		}
	}()

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)

	isLocalCall = err == nil

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
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(SayHelloResponse)
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
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// EchoSvc2Service
type EchoSvc2Service interface {
	entity.IService

	// SayHello
	SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _EchoSvc2Service_SayHello_Handler(svc any, ctx entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {
		return svc.(EchoSvc2Service).SayHello(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoSvc2Service/SayHello",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(EchoSvc2Service).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SayHelloResponse), err

}

func _EchoSvc2Service_SayHello_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _EchoSvc2Service_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), interceptor)

}

func _EchoSvc2Service_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _EchoSvc2Service_SayHello_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

var EchoSvc2ServiceDesc = entity.EntityDesc{
	TypeName:    "exampb.EchoSvc2Service",
	HandlerType: (*EchoSvc2Service)(nil),
	Methods: map[string]entity.MethodDesc{

		"SayHello": {
			MethodName:   "SayHello",
			Handler:      _EchoSvc2Service_SayHello_Remote_Handler,
			LocalHandler: _EchoSvc2Service_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}

func RegisterClientEntity(etyMgr *entity.EntityManager, impl ClientEntity) {
	etyMgr.RegisterEntityDesc(&ClientEntityDesc, impl)
}

func NewClientClient(eid entity.EntityID) *Client {
	t := &Client{
		EntityID: eid,
	}
	return t
}

func NewClient(c entity.Context) (*Client, error) {
	etyMgr := entity.GetEntityMgr(c)
	return NewClientWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewClientWithID(c entity.Context, id entity.EntityID) (*Client, error) {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, ClientEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil, err
	}
	t := &Client{
		EntityID: id,
	}

	return t, nil
}

type Client struct {
	EntityID entity.EntityID
}

func (t *Client) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// 如果是LocalCall，由local内部关闭chan
	isLocalCall := false
	var chanRsp chan any = make(chan any)
	defer func() {
		if !isLocalCall {
			close(chanRsp)
		}
	}()

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)

	isLocalCall = err == nil

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
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(SayHelloResponse)
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
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ClientEntity 客户端
type ClientEntity interface {
	entity.IEntity

	// SayHello
	SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _ClientEntity_SayHello_Handler(svc any, ctx entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {
		return svc.(ClientEntity).SayHello(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.ClientEntity/SayHello",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ClientEntity).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SayHelloResponse), err

}

func _ClientEntity_SayHello_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ClientEntity_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), interceptor)

}

func _ClientEntity_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _ClientEntity_SayHello_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

var ClientEntityDesc = entity.EntityDesc{
	TypeName:    "exampb.ClientEntity",
	HandlerType: (*ClientEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"SayHello": {
			MethodName:   "SayHello",
			Handler:      _ClientEntity_SayHello_Remote_Handler,
			LocalHandler: _ClientEntity_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}

func RegisterNoRespEntity(etyMgr *entity.EntityManager, impl NoRespEntity) {
	etyMgr.RegisterEntityDesc(&NoRespEntityDesc, impl)
}

func NewNoRespClient(eid entity.EntityID) *NoResp {
	t := &NoResp{
		EntityID: eid,
	}
	return t
}

func NewNoResp(c entity.Context) (*NoResp, error) {
	etyMgr := entity.GetEntityMgr(c)
	return NewNoRespWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewNoRespWithID(c entity.Context, id entity.EntityID) (*NoResp, error) {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, NoRespEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil, err
	}
	t := &NoResp{
		EntityID: id,
	}

	return t, nil
}

type NoResp struct {
	EntityID entity.EntityID
}

func (t *NoResp) SayHello(c entity.Context, req *SayHelloRequest) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_OneWay
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, nil)

	if err != nil {
		if !entity.NotMyMethod(err) {
			return err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return err
		}

		pkt.WriteBody(hdrBytes)

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

// NoRespEntity 客户端
type NoRespEntity interface {
	entity.IEntity

	// SayHello
	SayHello(c entity.Context, req *SayHelloRequest) error
}

func _NoRespEntity_SayHello_Handler(svc any, ctx entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(NoRespEntity).SayHello(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.NoRespEntity/SayHello",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(NoRespEntity).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _NoRespEntity_SayHello_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _NoRespEntity_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), interceptor)

}

func _NoRespEntity_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	if err != nil {
		// report
		return
	}
	err = _NoRespEntity_SayHello_Handler(svc, ctx, in, interceptor)
	_ = err
	// report err

}

var NoRespEntityDesc = entity.EntityDesc{
	TypeName:    "exampb.NoRespEntity",
	HandlerType: (*NoRespEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"SayHello": {
			MethodName:   "SayHello",
			Handler:      _NoRespEntity_SayHello_Remote_Handler,
			LocalHandler: _NoRespEntity_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}
