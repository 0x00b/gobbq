// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package exampb

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
	// exampb "github.com/0x00b/gobbq/example/exampb"
)

var _ = snowflake.GenUUID()

func RegisterEchoService(etyMgr *entity.EntityManager, impl EchoService) {
	etyMgr.RegisterService(&EchoServiceDesc, impl)
}

func NewEchoServiceClient(etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *echoService {
	t := &echoService{
		client: client,
		etyMgr: etyMgr,
	}
	return t
}

func NewEchoService(etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *echoService {
	t := &echoService{
		client: client,
		etyMgr: etyMgr,
	}
	return t
}

type echoService struct {
	etyMgr *entity.EntityManager
	client *codec.PacketReadWriter
}

func (t *echoService) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "exampb.EchoService"}
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := t.etyMgr.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

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

func _EchoService_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _EchoService_SayHello_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
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

func NewEchoEtyEntityClient(client *codec.PacketReadWriter, etyMgr *entity.EntityManager, entity *bbq.EntityID) *echoEtyEntity {
	t := &echoEtyEntity{
		client: client,
		etyMgr: etyMgr,
		entity: entity,
	}
	return t
}

func NewEchoEtyEntity(c entity.Context, etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *echoEtyEntity {
	return NewEchoEtyEntityWithID(c, etyMgr, etyMgr.EntityIDGenerator.NewEntityID("exampb.EchoEtyEntity"), client)
}

func NewEchoEtyEntityWithID(c entity.Context, etyMgr *entity.EntityManager, id *bbq.EntityID, client *codec.PacketReadWriter) *echoEtyEntity {

	_, err := etyMgr.NewEntity(c, id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &echoEtyEntity{
		entity: id,
		client: client,
		etyMgr: etyMgr,
	}

	return t
}

type echoEtyEntity struct {
	entity *bbq.EntityID

	etyMgr *entity.EntityManager
	client *codec.PacketReadWriter
}

func (t *echoEtyEntity) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := t.etyMgr.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

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

func _EchoEtyEntity_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _EchoEtyEntity_SayHello_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
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

func NewEchoSvc2ServiceClient(etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *echoSvc2Service {
	t := &echoSvc2Service{
		client: client,
		etyMgr: etyMgr,
	}
	return t
}

func NewEchoSvc2Service(etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *echoSvc2Service {
	t := &echoSvc2Service{
		client: client,
		etyMgr: etyMgr,
	}
	return t
}

type echoSvc2Service struct {
	etyMgr *entity.EntityManager
	client *codec.PacketReadWriter
}

func (t *echoSvc2Service) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "exampb.EchoSvc2Service"}
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := t.etyMgr.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

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

func _EchoSvc2Service_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _EchoSvc2Service_SayHello_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
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

func NewClientEntityClient(client *codec.PacketReadWriter, etyMgr *entity.EntityManager, entity *bbq.EntityID) *clientEntity {
	t := &clientEntity{
		client: client,
		etyMgr: etyMgr,
		entity: entity,
	}
	return t
}

func NewClientEntity(c entity.Context, etyMgr *entity.EntityManager, client *codec.PacketReadWriter) *clientEntity {
	return NewClientEntityWithID(c, etyMgr, etyMgr.EntityIDGenerator.NewEntityID("exampb.ClientEntity"), client)
}

func NewClientEntityWithID(c entity.Context, etyMgr *entity.EntityManager, id *bbq.EntityID, client *codec.PacketReadWriter) *clientEntity {

	_, err := etyMgr.NewEntity(c, id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &clientEntity{
		entity: id,
		client: client,
		etyMgr: etyMgr,
	}

	return t
}

type clientEntity struct {
	entity *bbq.EntityID

	etyMgr *entity.EntityManager
	client *codec.PacketReadWriter
}

func (t *clientEntity) SayHello(c entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := t.etyMgr.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

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

func _ClientEntity_SayHello_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ClientEntity_SayHello_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
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
