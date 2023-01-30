// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package proxypb

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	// proxypb "github.com/0x00b/gobbq/components/proxy/proxypb"
)

func RegisterProxyService(impl ProxyService) {
	entity.Manager.RegisterService(&ProxyServiceDesc, impl)
}

func NewProxyService() *proxyService {
	t := &proxyService{}
	return t
}

type proxyService struct {
}

func (t *proxyService) RegisterEntity(c context.Context, req *RegisterEntityRequest, callback func(c context.Context, rsp *RegisterEntityResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.Method = "proxypb.ProxyService/RegisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterEntityResponse))
	}

	err = entity.HandleCallLocalMethod(c, pkt, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)

		//todo get response
		var requestMap map[string]func(c context.Context, rsp interface{})
		requestMap[pkt.Header.RequestId] = itfCallback

		if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
			cb := requestMap[pkt.Header.RequestId]

			cb(context.Background(), nil)

		}

	}

	return err

}

func (t *proxyService) UnregisterEntity(c context.Context, req *RegisterEntityRequest, callback func(c context.Context, rsp *RegisterEntityResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.Method = "proxypb.ProxyService/UnregisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterEntityResponse))
	}

	err = entity.HandleCallLocalMethod(c, pkt, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)

		//todo get response
		var requestMap map[string]func(c context.Context, rsp interface{})
		requestMap[pkt.Header.RequestId] = itfCallback

		if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
			cb := requestMap[pkt.Header.RequestId]

			cb(context.Background(), nil)

		}

	}

	return err

}

func (t *proxyService) Ping(c context.Context, req *PingPong, callback func(c context.Context, rsp *PingPong)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = nil
	pkt.Header.Method = "proxypb.ProxyService/Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*PingPong))
	}

	err = entity.HandleCallLocalMethod(c, pkt, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)

		//todo get response
		var requestMap map[string]func(c context.Context, rsp interface{})
		requestMap[pkt.Header.RequestId] = itfCallback

		if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
			cb := requestMap[pkt.Header.RequestId]

			cb(context.Background(), nil)

		}

	}

	return err

}

// ProxyService
type ProxyService interface {
	entity.IService

	// RegisterEntity
	RegisterEntity(c context.Context, req *RegisterEntityRequest, ret func(*RegisterEntityResponse, error))

	// UnregisterEntity
	UnregisterEntity(c context.Context, req *RegisterEntityRequest, ret func(*RegisterEntityResponse, error))

	// Ping
	Ping(c context.Context, req *PingPong, ret func(*PingPong, error))
}

func _ProxyService_RegisterEntity_Handler(svc interface{}, ctx context.Context, in *RegisterEntityRequest, ret func(rsp *RegisterEntityResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).RegisterEntity(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterEntity",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(ProxyService).RegisterEntity(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*RegisterEntityResponse), err) }, handler)
	return
}

func _ProxyService_RegisterEntity_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterEntityResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_RegisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), ret, interceptor)
	return
}

func _ProxyService_RegisterEntity_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterEntityResponse, err error) {

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version = hdr.Version
		npkt.Header.RequestId = hdr.RequestId
		npkt.Header.Timeout = hdr.Timeout
		npkt.Header.RequestType = hdr.RequestType
		npkt.Header.ServiceType = hdr.ServiceType
		npkt.Header.SrcEntity = hdr.DstEntity
		npkt.Header.DstEntity = hdr.SrcEntity
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

func _ProxyService_UnregisterEntity_Handler(svc interface{}, ctx context.Context, in *RegisterEntityRequest, ret func(rsp *RegisterEntityResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).UnregisterEntity(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/UnregisterEntity",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(ProxyService).UnregisterEntity(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*RegisterEntityResponse), err) }, handler)
	return
}

func _ProxyService_UnregisterEntity_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {

	ret := func(rsp *RegisterEntityResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_UnregisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), ret, interceptor)
	return
}

func _ProxyService_UnregisterEntity_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *RegisterEntityResponse, err error) {

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version = hdr.Version
		npkt.Header.RequestId = hdr.RequestId
		npkt.Header.Timeout = hdr.Timeout
		npkt.Header.RequestType = hdr.RequestType
		npkt.Header.ServiceType = hdr.ServiceType
		npkt.Header.SrcEntity = hdr.DstEntity
		npkt.Header.DstEntity = hdr.SrcEntity
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

func _ProxyService_Ping_Handler(svc interface{}, ctx context.Context, in *PingPong, ret func(rsp *PingPong, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(ProxyService).Ping(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/Ping",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(ProxyService).Ping(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*PingPong), err) }, handler)
	return
}

func _ProxyService_Ping_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {

	ret := func(rsp *PingPong, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}

	_ProxyService_Ping_Handler(svc, ctx, in.(*PingPong), ret, interceptor)
	return
}

func _ProxyService_Ping_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	ret := func(rsp *PingPong, err error) {

		npkt, release := codec.NewPacket()
		defer release()

		npkt.Header.Version = hdr.Version
		npkt.Header.RequestId = hdr.RequestId
		npkt.Header.Timeout = hdr.Timeout
		npkt.Header.RequestType = hdr.RequestType
		npkt.Header.ServiceType = hdr.ServiceType
		npkt.Header.SrcEntity = hdr.DstEntity
		npkt.Header.DstEntity = hdr.SrcEntity
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

		"UnregisterEntity": {
			MethodName:   "UnregisterEntity",
			Handler:      _ProxyService_UnregisterEntity_Remote_Handler,
			LocalHandler: _ProxyService_UnregisterEntity_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxyService_Ping_Remote_Handler,
			LocalHandler: _ProxyService_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}
