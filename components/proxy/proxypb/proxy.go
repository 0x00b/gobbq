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

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Service,
		SrcEntity:    nil,
		DstEntity:    nil,
		Method:       "proxypb.ProxyService/RegisterEntity",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterEntityResponse))
	}

	err = entity.HandleCallLocalMethod(c, hdr, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		pkt := codec.NewPacket()

		pkt.SetHeader(hdr)

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)
		//todo get response

	}

	return err

}

func (t *proxyService) UnregisterEntity(c context.Context, req *RegisterEntityRequest, callback func(c context.Context, rsp *RegisterEntityResponse)) (err error) {

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Service,
		SrcEntity:    nil,
		DstEntity:    nil,
		Method:       "proxypb.ProxyService/UnregisterEntity",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterEntityResponse))
	}

	err = entity.HandleCallLocalMethod(c, hdr, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		pkt := codec.NewPacket()

		pkt.SetHeader(hdr)

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)
		//todo get response

	}

	return err

}

func (t *proxyService) Ping(c context.Context, req *PingPong, callback func(c context.Context, rsp *PingPong)) (err error) {

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Service,
		SrcEntity:    nil,
		DstEntity:    nil,
		Method:       "proxypb.ProxyService/Ping",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*PingPong))
	}

	err = entity.HandleCallLocalMethod(c, hdr, req, itfCallback)
	if err == nil {
		return nil
	}

	if entity.NotMyMethod(err) {

		pkt := codec.NewPacket()

		pkt.SetHeader(hdr)

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			return err
		}

		pkt.WriteBody(hdrBytes)

		ex.SendProxy(pkt)
		//todo get response

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

	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterEntityRequest)

	ret := func(rsp *RegisterEntityResponse, err error) {
		npkt := codec.NewPacket()

		rhdr := &bbq.Header{
			Version:      hdr.Version,
			RequestId:    hdr.RequestId,
			Timeout:      hdr.Timeout,
			RequestType:  hdr.RequestType,
			ServiceType:  hdr.ServiceType,
			SrcEntity:    hdr.DstEntity,
			DstEntity:    hdr.SrcEntity,
			Method:       hdr.Method,
			ContentType:  hdr.ContentType,
			CompressType: hdr.CompressType,
			CheckFlags:   0,
			TransInfo:    hdr.TransInfo,
			ErrCode:      0,
			ErrMsg:       "",
		}
		npkt.SetHeader(rhdr)

		rbyte, err := codec.DefaultCodec.Marshal(rhdr)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
		npkt.WriteBody(rbyte)

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

	if err := dec(in); err != nil {
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

	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterEntityRequest)

	ret := func(rsp *RegisterEntityResponse, err error) {
		npkt := codec.NewPacket()

		rhdr := &bbq.Header{
			Version:      hdr.Version,
			RequestId:    hdr.RequestId,
			Timeout:      hdr.Timeout,
			RequestType:  hdr.RequestType,
			ServiceType:  hdr.ServiceType,
			SrcEntity:    hdr.DstEntity,
			DstEntity:    hdr.SrcEntity,
			Method:       hdr.Method,
			ContentType:  hdr.ContentType,
			CompressType: hdr.CompressType,
			CheckFlags:   0,
			TransInfo:    hdr.TransInfo,
			ErrCode:      0,
			ErrMsg:       "",
		}
		npkt.SetHeader(rhdr)

		rbyte, err := codec.DefaultCodec.Marshal(rhdr)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
		npkt.WriteBody(rbyte)

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

	if err := dec(in); err != nil {
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

	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(PingPong)

	ret := func(rsp *PingPong, err error) {
		npkt := codec.NewPacket()

		rhdr := &bbq.Header{
			Version:      hdr.Version,
			RequestId:    hdr.RequestId,
			Timeout:      hdr.Timeout,
			RequestType:  hdr.RequestType,
			ServiceType:  hdr.ServiceType,
			SrcEntity:    hdr.DstEntity,
			DstEntity:    hdr.SrcEntity,
			Method:       hdr.Method,
			ContentType:  hdr.ContentType,
			CompressType: hdr.CompressType,
			CheckFlags:   0,
			TransInfo:    hdr.TransInfo,
			ErrCode:      0,
			ErrMsg:       "",
		}
		npkt.SetHeader(rhdr)

		rbyte, err := codec.DefaultCodec.Marshal(rhdr)
		if err != nil {
			fmt.Println("WritePacket", err)
			return
		}
		npkt.WriteBody(rbyte)

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

	if err := dec(in); err != nil {
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
