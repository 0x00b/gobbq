// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package gatepb

import (
	"context"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/proto/bbq"
	"fmt"

	// gatepb "github.com/0x00b/gobbq/components/gate/gatepb"

)

func RegisterGateEntity(impl GateEntity) {
	entity.Manager.RegisterEntity(&GateEntityDesc, impl)
}

func NewGateEntity() *gateEntity {
	return NewGateEntityWithID(entity.EntityID(snowflake.GenUUID()))
}

func NewGateEntityWithID(id entity.EntityID) *gateEntity {

	ety := entity.NewEntity(id, GateEntityDesc.TypeName)
	t := &gateEntity{entity: ety}

	return t
}

type gateEntity struct {
	entity *bbq.EntityID
}

func (t *gateEntity) RegisterClient(c context.Context, req *RegisterClientRequest, callback func(c context.Context, rsp *RegisterClientResponse)) (err error) {

	pkt := codec.NewPacket()

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Entity,
		SrcEntity:    nil,
		DstEntity:    t.entity,
		Method:       "gatepb.GateEntity/RegisterClient",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil

}

func (t *gateEntity) UnregisterClient(c context.Context, req *RegisterClientRequest, callback func(c context.Context, rsp *RegisterClientResponse)) (err error) {

	pkt := codec.NewPacket()

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Entity,
		SrcEntity:    nil,
		DstEntity:    t.entity,
		Method:       "gatepb.GateEntity/UnregisterClient",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil

}

func (t *gateEntity) Ping(c context.Context, req *PingPong, callback func(c context.Context, rsp *PingPong)) (err error) {

	pkt := codec.NewPacket()

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  bbq.RequestType_RequestRequest,
		ServiceType:  bbq.ServiceType_Entity,
		SrcEntity:    nil,
		DstEntity:    t.entity,
		Method:       "gatepb.GateEntity/Ping",
		ContentType:  bbq.ContentType_Proto,
		CompressType: bbq.CompressType_None,
		CheckFlags:   0,
		TransInfo:    map[string][]byte{},
		ErrCode:      0,
		ErrMsg:       "",
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		return err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil

}

// GateEntity
type GateEntity interface {
	entity.IEntity

	// RegisterClient
	RegisterClient(c context.Context, req *RegisterClientRequest, ret func(*RegisterClientResponse, error))

	// UnregisterClient
	UnregisterClient(c context.Context, req *RegisterClientRequest, ret func(*RegisterClientResponse, error))

	// Ping
	Ping(c context.Context, req *PingPong, ret func(*PingPong, error))
}

func _GateEntity_RegisterClient_Handler(svc interface{}, ctx context.Context, in *RegisterClientRequest, ret func(rsp *RegisterClientResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateEntity).RegisterClient(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateEntity/RegisterClient",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(GateEntity).RegisterClient(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*RegisterClientResponse), err) }, handler)
	return
}

func _GateEntity_RegisterClient_Local_Handler(svc interface{}, ctx context.Context, in *RegisterClientRequest, callback func(c context.Context, rsp *RegisterClientResponse), interceptor entity.ServerInterceptor) {
	ret := func(rsp *RegisterClientResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_RegisterClient_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateEntity_RegisterClient_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterClientRequest)

	ret := func(rsp *RegisterClientResponse, err error) {
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

	_GateEntity_RegisterClient_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateEntity_UnregisterClient_Handler(svc interface{}, ctx context.Context, in *RegisterClientRequest, ret func(rsp *RegisterClientResponse, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateEntity).UnregisterClient(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateEntity/UnregisterClient",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(GateEntity).UnregisterClient(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*RegisterClientResponse), err) }, handler)
	return
}

func _GateEntity_UnregisterClient_Local_Handler(svc interface{}, ctx context.Context, in *RegisterClientRequest, callback func(c context.Context, rsp *RegisterClientResponse), interceptor entity.ServerInterceptor) {
	ret := func(rsp *RegisterClientResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_UnregisterClient_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateEntity_UnregisterClient_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.GetHeader()
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterClientRequest)

	ret := func(rsp *RegisterClientResponse, err error) {
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

	_GateEntity_UnregisterClient_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateEntity_Ping_Handler(svc interface{}, ctx context.Context, in *PingPong, ret func(rsp *PingPong, err error), interceptor entity.ServerInterceptor) {
	if interceptor == nil {
		svc.(GateEntity).Ping(ctx, in, ret)
		return
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateEntity/Ping",
	}

	handler := func(ctx context.Context, rsp interface{}, _ entity.RetFunc) {
		svc.(GateEntity).Ping(ctx, in, ret)
	}

	interceptor(ctx, in, info, func(i interface{}, err error) { ret(i.(*PingPong), err) }, handler)
	return
}

func _GateEntity_Ping_Local_Handler(svc interface{}, ctx context.Context, in *PingPong, callback func(c context.Context, rsp *PingPong), interceptor entity.ServerInterceptor) {
	ret := func(rsp *PingPong, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_Ping_Handler(svc, ctx, in, ret, interceptor)
	return
}

func _GateEntity_Ping_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

	_GateEntity_Ping_Handler(svc, ctx, in, ret, interceptor)
	return
}

var GateEntityDesc = entity.ServiceDesc{
	TypeName:    "gatepb.GateEntity",
	HandlerType: (*GateEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterClient": {
			MethodName: "RegisterClient",
			Handler:    _GateEntity_RegisterClient_Remote_Handler,
		},

		"UnregisterClient": {
			MethodName: "UnregisterClient",
			Handler:    _GateEntity_UnregisterClient_Remote_Handler,
		},

		"Ping": {
			MethodName: "Ping",
			Handler:    _GateEntity_Ping_Remote_Handler,
		},
	},

	Metadata: "gate.proto",
}
