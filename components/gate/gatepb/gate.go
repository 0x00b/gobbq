// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package gatepb

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"fmt"

	// gatepb "github.com/0x00b/gobbq/components/gate/gatepb"

)

var _ = snowflake.GenUUID()

func RegisterGateService(impl GateService) {
	entity.Manager.RegisterService(&GateServiceDesc, impl)
}

func NewGateServiceClient(client *nets.Client) *gateService {
	t := &gateService{client: client}
	return t
}

func NewGateService(client *nets.Client) *gateService {
	t := &gateService{client: client}
	return t
}

type gateService struct {
	client *nets.Client
}

func (t *gateService) RegisterClient(c *entity.Context, req *RegisterClientRequest, callback func(c *entity.Context, rsp *RegisterClientResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "gatepb.GateService"
	pkt.Header.Method = "RegisterClient"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*RegisterClientResponse))
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

func (t *gateService) UnregisterClient(c *entity.Context, req *RegisterClientRequest) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "gatepb.GateService"
	pkt.Header.Method = "UnregisterClient"
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

	// if entity.NotMyMethod(err) {

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// }

	return err

}

func (t *gateService) Ping(c *entity.Context, req *PingPong, callback func(c *entity.Context, rsp *PingPong)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "gatepb.GateService"
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*PingPong))
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

// GateService
type GateService interface {
	entity.IService

	// RegisterClient
	RegisterClient(c *entity.Context, req *RegisterClientRequest, ret func(*RegisterClientResponse, error))

	// UnregisterClient
	UnregisterClient(c *entity.Context, req *RegisterClientRequest)

	// Ping
	Ping(c *entity.Context, req *PingPong, ret func(*PingPong, error))
}

func _GateService_RegisterClient_Handler(svc any, ctx *entity.Context, in *RegisterClientRequest, ret func(rsp *RegisterClientResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateService).RegisterClient(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/RegisterClient",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(GateService).RegisterClient(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*RegisterClientResponse), err) }, handler)
	return
}

func _GateService_RegisterClient_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterClientResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_GateService_RegisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), ret, interceptor)
	return
}

func _GateService_RegisterClient_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterClientResponse, err error) {

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

	in := new(RegisterClientRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_GateService_RegisterClient_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateService_UnregisterClient_Handler(svc any, ctx *entity.Context, in *RegisterClientRequest, interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateService).UnregisterClient(ctx, in)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/UnregisterClient",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(GateService).UnregisterClient(ctx, in)
	}

	interceptor(ctx, in, info, nil, handler)
	return
}

func _GateService_UnregisterClient_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	_GateService_UnregisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), interceptor)
	return
}

func _GateService_UnregisterClient_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterClientRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		return
	}

	_GateService_UnregisterClient_Handler(svc, ctx, in, interceptor)
	return
}

func _GateService_Ping_Handler(svc any, ctx *entity.Context, in *PingPong, ret func(rsp *PingPong, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateService).Ping(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/Ping",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(GateService).Ping(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*PingPong), err) }, handler)
	return
}

func _GateService_Ping_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *PingPong, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_GateService_Ping_Handler(svc, ctx, in.(*PingPong), ret, interceptor)
	return
}

func _GateService_Ping_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *PingPong, err error) {

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

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_GateService_Ping_Handler(svc, ctx, in, ret, interceptor)
	return
}

var GateServiceDesc = entity.ServiceDesc{
	TypeName:    "gatepb.GateService",
	HandlerType: (*GateService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterClient": {
			MethodName:   "RegisterClient",
			Handler:      _GateService_RegisterClient_Remote_Handler,
			LocalHandler: _GateService_RegisterClient_Local_Handler,
		},

		"UnregisterClient": {
			MethodName:   "UnregisterClient",
			Handler:      _GateService_UnregisterClient_Remote_Handler,
			LocalHandler: _GateService_UnregisterClient_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _GateService_Ping_Remote_Handler,
			LocalHandler: _GateService_Ping_Local_Handler,
		},
	},

	Metadata: "gate.proto",
}
