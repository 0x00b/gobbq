// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package exampb

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
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

func (t *echoService) SayHello(c *entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	eid := ""
	if c != nil {
		eid = string(c.EntityID())
	}
	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = eid
	pkt.Header.DstEntity = ""
	pkt.Header.ServiceName = "exampb.EchoService"
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.Entity.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// EchoService
type EchoService interface {
	entity.IEntity

	// SayHello
	SayHello(c *entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _EchoService_SayHello_Handler(svc any, ctx *entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {

		return svc.(EchoService).SayHello(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoService/SayHello",
	}

	handler := func(ctx *entity.Context, rsp any) (any, error) {

		return svc.(EchoService).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*SayHelloResponse), err

}

//func _EchoService_SayHello_Local_Handler(svc any, ctx *entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *SayHelloResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_EchoService_SayHello_Handler(svc, ctx, in.(*SayHelloRequest) , ret, interceptor)
//
//}

func _EchoService_SayHello_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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
	npkt.Header.ServiceName = hdr.ServiceName
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
			MethodName: "SayHello",
			Handler:    _EchoService_SayHello_Remote_Handler,
			//LocalHandler:	_EchoService_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}

func RegisterEchoEtyEntity(impl EchoEtyEntity) {
	entity.Manager.RegisterEntity(&EchoEtyEntityDesc, impl)
}

func NewEchoEtyEntityClient(client *nets.Client, entity entity.EntityID) *echoEtyEntity {
	t := &echoEtyEntity{client: client, entity: entity}
	return t
}

func NewEchoEtyEntity(c *entity.Context, client *nets.Client) *echoEtyEntity {
	return NewEchoEtyEntityWithID(c, entity.EntityID(snowflake.GenUUID()), client)
}

func NewEchoEtyEntityWithID(c *entity.Context, id entity.EntityID, client *nets.Client) *echoEtyEntity {

	err := entity.NewEntity(c, &id, EchoEtyEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &echoEtyEntity{entity: id, client: client}

	return t
}

type echoEtyEntity struct {
	entity entity.EntityID

	client *nets.Client
}

func (t *echoEtyEntity) SayHello(c *entity.Context, req *SayHelloRequest) (*SayHelloResponse, error) {

	eid := ""
	if c != nil {
		eid = string(c.EntityID())
	}
	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = eid
	pkt.Header.DstEntity = string(t.entity)
	pkt.Header.ServiceName = "exampb.EchoEtyEntity"
	pkt.Header.Method = "SayHello"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.Entity.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SayHelloResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*SayHelloResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// EchoEtyEntity
type EchoEtyEntity interface {
	entity.IEntity

	// SayHello
	SayHello(c *entity.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

func _EchoEtyEntity_SayHello_Handler(svc any, ctx *entity.Context, in *SayHelloRequest, interceptor entity.ServerInterceptor) (*SayHelloResponse, error) {
	if interceptor == nil {

		return svc.(EchoEtyEntity).SayHello(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/exampb.EchoEtyEntity/SayHello",
	}

	handler := func(ctx *entity.Context, rsp any) (any, error) {

		return svc.(EchoEtyEntity).SayHello(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*SayHelloResponse), err

}

//func _EchoEtyEntity_SayHello_Local_Handler(svc any, ctx *entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *SayHelloResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_EchoEtyEntity_SayHello_Handler(svc, ctx, in.(*SayHelloRequest) , ret, interceptor)
//
//}

func _EchoEtyEntity_SayHello_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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
	npkt.Header.ServiceName = hdr.ServiceName
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
			MethodName: "SayHello",
			Handler:    _EchoEtyEntity_SayHello_Remote_Handler,
			//LocalHandler:	_EchoEtyEntity_SayHello_Local_Handler,
		},
	},

	Metadata: "exam.proto",
}
