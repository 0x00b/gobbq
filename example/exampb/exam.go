// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package exampb

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"fmt"

	// exampb "github.com/0x00b/gobbq/example/exampb"

)

var _ = snowflake.GenUUID()

func RegisterEchoService(impl EchoService) {
	entity.Manager.RegisterService(&EchoServiceDesc, impl)
}

func NewEchoServiceClient(client *nets.Client) *echoService {
	t := &echoService{client: client}
	return t
}

func NewEchoService(client *nets.Client) *echoService {
	t := &echoService{client: client}
	return t
}

type echoService struct {
	client *nets.Client
}

func (t *echoService) SayHello(c *entity.Context, req *SayHelloRequest, callback func(c *entity.Context, rsp *SayHelloResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "exampb.EchoService"
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*SayHelloResponse))
	}
	_ = itfCallback

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	// if entity.NotMyMethod(err) {

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback for request
	// c.Service.RegisterCallback(pkt.Header.RequestId, itfCallback)

	// }

	return err

}

// EchoService
type EchoService interface {
	entity.IService

	// SayHello
	SayHello(c *entity.Context, req *SayHelloRequest, ret func(*SayHelloResponse, error))
}

func _EchoService_SayHello_Handler(svc any, ctx *entity.Context, in *SayHelloRequest, ret func(rsp *SayHelloResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(EchoService).SayHello(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoService/SayHello",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(EchoService).SayHello(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*SayHelloResponse), err) }, handler)
	return
}

func _EchoService_SayHello_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *SayHelloResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_EchoService_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), ret, interceptor)
	return
}

func _EchoService_SayHello_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *SayHelloResponse, err error) {

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version = hdr.Version
		npkt.Header.RequestId = hdr.RequestId
		npkt.Header.Timeout = hdr.Timeout
		npkt.Header.RequestType = bbq.RequestType_RequestRespone
		npkt.Header.ServiceType = hdr.ServiceType
		npkt.Header.SrcEntity = hdr.DstEntity
		npkt.Header.DstEntity = hdr.SrcEntity
		npkt.Header.ServiceName = hdr.ServiceName
		npkt.Header.Method = hdr.Method
		npkt.Header.ContentType = hdr.ContentType
		npkt.Header.CompressType = hdr.CompressType
		npkt.Header.CheckFlags = 0
		npkt.Header.TransInfo = hdr.TransInfo
		npkt.Header.ErrCode = 0
		npkt.Header.ErrMsg = ""

		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			fmt.Println("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)

		err = pkt.Src.WritePacket(npkt)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
	}

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_EchoService_SayHello_Handler(svc, ctx, in, ret, interceptor)
	return
}

var EchoServiceDesc = entity.ServiceDesc{
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

func RegisterEchoEtyEntity(impl EchoEtyEntity) {
	entity.Manager.RegisterEntity(&EchoEtyEntityDesc, impl)
}

func NewEchoEtyEntityClient(client *nets.Client, entity *bbq.EntityID) *echoEtyEntity {
	t := &echoEtyEntity{client: client, entity: entity}
	return t
}

func NewEchoEtyEntity(c *entity.Context, client *nets.Client) *echoEtyEntity {
	return NewEchoEtyEntityWithID(c, entity.EntityID(snowflake.GenUUID()), client)
}

func NewEchoEtyEntityWithID(c *entity.Context, id entity.EntityID, client *nets.Client) *echoEtyEntity {

	ety := entity.NewEntity(c, id, EchoEtyEntityDesc.TypeName)
	t := &echoEtyEntity{entity: ety, client: client}

	return t
}

type echoEtyEntity struct {
	entity *bbq.EntityID

	client *nets.Client
}

func (t *echoEtyEntity) SayHello(c *entity.Context, req *SayHelloRequest, callback func(c *entity.Context, rsp *SayHelloResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = t.entity
	pkt.Header.ServiceName = "exampb.EchoEtyEntity"
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*SayHelloResponse))
	}
	_ = itfCallback

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	// if entity.NotMyMethod(err) {

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback for request
	// c.Service.RegisterCallback(pkt.Header.RequestId, itfCallback)

	// }

	return err

}

// EchoEtyEntity
type EchoEtyEntity interface {
	entity.IEntity

	// SayHello
	SayHello(c *entity.Context, req *SayHelloRequest, ret func(*SayHelloResponse, error))
}

func _EchoEtyEntity_SayHello_Handler(svc any, ctx *entity.Context, in *SayHelloRequest, ret func(rsp *SayHelloResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(EchoEtyEntity).SayHello(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoEtyEntity/SayHello",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(EchoEtyEntity).SayHello(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*SayHelloResponse), err) }, handler)
	return
}

func _EchoEtyEntity_SayHello_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *SayHelloResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_EchoEtyEntity_SayHello_Handler(svc, ctx, in.(*SayHelloRequest), ret, interceptor)
	return
}

func _EchoEtyEntity_SayHello_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *SayHelloResponse, err error) {

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version = hdr.Version
		npkt.Header.RequestId = hdr.RequestId
		npkt.Header.Timeout = hdr.Timeout
		npkt.Header.RequestType = bbq.RequestType_RequestRespone
		npkt.Header.ServiceType = hdr.ServiceType
		npkt.Header.SrcEntity = hdr.DstEntity
		npkt.Header.DstEntity = hdr.SrcEntity
		npkt.Header.ServiceName = hdr.ServiceName
		npkt.Header.Method = hdr.Method
		npkt.Header.ContentType = hdr.ContentType
		npkt.Header.CompressType = hdr.CompressType
		npkt.Header.CheckFlags = 0
		npkt.Header.TransInfo = hdr.TransInfo
		npkt.Header.ErrCode = 0
		npkt.Header.ErrMsg = ""

		rb, err := codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			fmt.Println("Marshal(rsp)", err)
			return
		}

		npkt.WriteBody(rb)

		err = pkt.Src.WritePacket(npkt)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
	}

	in := new(SayHelloRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_EchoEtyEntity_SayHello_Handler(svc, ctx, in, ret, interceptor)
	return
}

var EchoEtyEntityDesc = entity.ServiceDesc{
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
