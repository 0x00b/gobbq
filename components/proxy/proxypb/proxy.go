// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package proxypb

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"fmt"

	// proxypb "github.com/0x00b/gobbq/components/proxy/proxypb"

)

var _ = snowflake.GenUUID()

func RegisterProxyService(impl ProxyService) {
	entity.Manager.RegisterService(&ProxyServiceDesc, impl)
}

func NewProxyServiceClient(client *nets.Client) *proxyService {
	t := &proxyService{client: client}
	return t
}

func NewProxyService(client *nets.Client) *proxyService {
	t := &proxyService{client: client}
	return t
}

type proxyService struct {
	client *nets.Client
}

func (t *proxyService) RegisterEntity(c *entity.Context, req *RegisterEntityRequest, callback func(c *entity.Context, rsp *RegisterEntityResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "proxypb.ProxyService"
	pkt.Header.Method = "RegisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*RegisterEntityResponse))
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

func (t *proxyService) RegisterService(c *entity.Context, req *RegisterServiceRequest, callback func(c *entity.Context, rsp *RegisterServiceResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "proxypb.ProxyService"
	pkt.Header.Method = "RegisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*RegisterServiceResponse))
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

func (t *proxyService) UnregisterEntity(c *entity.Context, req *RegisterEntityRequest, callback func(c *entity.Context, rsp *RegisterEntityResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "proxypb.ProxyService"
	pkt.Header.Method = "UnregisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*RegisterEntityResponse))
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

func (t *proxyService) UnregisterService(c *entity.Context, req *RegisterServiceRequest, callback func(c *entity.Context, rsp *RegisterServiceResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "proxypb.ProxyService"
	pkt.Header.Method = "UnregisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c *entity.Context, rsp any) {
		callback(c, rsp.(*RegisterServiceResponse))
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

func (t *proxyService) Ping(c *entity.Context, req *PingPong, callback func(c *entity.Context, rsp *PingPong)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.ServiceName = "proxypb.ProxyService"
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

// ProxyService
type ProxyService interface {
	entity.IService

	// RegisterEntity
	RegisterEntity(c *entity.Context, req *RegisterEntityRequest, ret func(*RegisterEntityResponse, error))

	// RegisterService
	RegisterService(c *entity.Context, req *RegisterServiceRequest, ret func(*RegisterServiceResponse, error))

	// UnregisterEntity
	UnregisterEntity(c *entity.Context, req *RegisterEntityRequest, ret func(*RegisterEntityResponse, error))

	// UnregisterService
	UnregisterService(c *entity.Context, req *RegisterServiceRequest, ret func(*RegisterServiceResponse, error))

	// Ping
	Ping(c *entity.Context, req *PingPong, ret func(*PingPong, error))
}

func _ProxyService_RegisterEntity_Handler(svc any, ctx *entity.Context, in *RegisterEntityRequest, ret func(rsp *RegisterEntityResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).RegisterEntity(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterEntity",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(ProxyService).RegisterEntity(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*RegisterEntityResponse), err) }, handler)
	return
}

func _ProxyService_RegisterEntity_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterEntityResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_RegisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), ret, interceptor)
	return
}

func _ProxyService_RegisterEntity_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterEntityResponse, err error) {

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

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_ProxyService_RegisterEntity_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _ProxyService_RegisterService_Handler(svc any, ctx *entity.Context, in *RegisterServiceRequest, ret func(rsp *RegisterServiceResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).RegisterService(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterService",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(ProxyService).RegisterService(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*RegisterServiceResponse), err) }, handler)
	return
}

func _ProxyService_RegisterService_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterServiceResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_RegisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), ret, interceptor)
	return
}

func _ProxyService_RegisterService_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterServiceResponse, err error) {

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

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_ProxyService_RegisterService_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _ProxyService_UnregisterEntity_Handler(svc any, ctx *entity.Context, in *RegisterEntityRequest, ret func(rsp *RegisterEntityResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).UnregisterEntity(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/UnregisterEntity",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(ProxyService).UnregisterEntity(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*RegisterEntityResponse), err) }, handler)
	return
}

func _ProxyService_UnregisterEntity_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterEntityResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_UnregisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), ret, interceptor)
	return
}

func _ProxyService_UnregisterEntity_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterEntityResponse, err error) {

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

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_ProxyService_UnregisterEntity_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _ProxyService_UnregisterService_Handler(svc any, ctx *entity.Context, in *RegisterServiceRequest, ret func(rsp *RegisterServiceResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).UnregisterService(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/UnregisterService",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(ProxyService).UnregisterService(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*RegisterServiceResponse), err) }, handler)
	return
}

func _ProxyService_UnregisterService_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterServiceResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_UnregisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), ret, interceptor)
	return
}

func _ProxyService_UnregisterService_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterServiceResponse, err error) {

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

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		ret(nil, err)
		return
	}

	_ProxyService_UnregisterService_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _ProxyService_Ping_Handler(svc any, ctx *entity.Context, in *PingPong, ret func(rsp *PingPong, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).Ping(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/Ping",
	}

	handler := func(ctx *entity.Context, rsp any, _ entity.RetFunc) {
		svc.(ProxyService).Ping(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i any, err error) { ret(i.(*PingPong), err) }, handler)
	return
}

func _ProxyService_Ping_Local_Handler(svc any, ctx *entity.Context, in any, callback func(c *entity.Context, rsp any), interceptor entity.ServerInterceptor) {

	ret := func(rsp *PingPong, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_Ping_Handler(svc, ctx, in.(*PingPong), ret, interceptor)
	return
}

func _ProxyService_Ping_Remote_Handler(svc any, ctx *entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

	_ProxyService_Ping_Handler(svc, ctx, in, ret, interceptor)
	return
}

var ProxyServiceDesc = entity.ServiceDesc{
	TypeName:    "proxypb.ProxyService",
	HandlerType: (*ProxyService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterEntity": {
			MethodName:   "RegisterEntity",
			Handler:      _ProxyService_RegisterEntity_Remote_Handler,
			LocalHandler: _ProxyService_RegisterEntity_Local_Handler,
		},

		"RegisterService": {
			MethodName:   "RegisterService",
			Handler:      _ProxyService_RegisterService_Remote_Handler,
			LocalHandler: _ProxyService_RegisterService_Local_Handler,
		},

		"UnregisterEntity": {
			MethodName:   "UnregisterEntity",
			Handler:      _ProxyService_UnregisterEntity_Remote_Handler,
			LocalHandler: _ProxyService_UnregisterEntity_Local_Handler,
		},

		"UnregisterService": {
			MethodName:   "UnregisterService",
			Handler:      _ProxyService_UnregisterService_Remote_Handler,
			LocalHandler: _ProxyService_UnregisterService_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxyService_Ping_Remote_Handler,
			LocalHandler: _ProxyService_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}
