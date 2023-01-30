// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package gatepb

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
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

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "gatepb.GateEntity/RegisterClient"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterClientResponse))
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

	}

	return err

}

func (t *gateEntity) UnregisterClient(c context.Context, req *RegisterClientRequest, callback func(c context.Context, rsp *RegisterClientResponse)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "gatepb.GateEntity/UnregisterClient"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	itfCallback := func(c context.Context, rsp interface{}) {
		callback(c, rsp.(*RegisterClientResponse))
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

	}

	return err

}

func (t *gateEntity) Ping(c context.Context, req *PingPong, callback func(c context.Context, rsp *PingPong)) (err error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = nil
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "gatepb.GateEntity/Ping"
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

	}

	return err

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

func _GateEntity_RegisterClient_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {
	ret := func(rsp *RegisterClientResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_RegisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), ret, interceptor)
	return
}

func _GateEntity_RegisterClient_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterClientRequest)

	ret := func(rsp *RegisterClientResponse, err error) {

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

func _GateEntity_UnregisterClient_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {
	ret := func(rsp *RegisterClientResponse, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_UnregisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), ret, interceptor)
	return
}

func _GateEntity_UnregisterClient_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(RegisterClientRequest)

	ret := func(rsp *RegisterClientResponse, err error) {

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

func _GateEntity_Ping_Local_Handler(svc interface{}, ctx context.Context, in interface{}, callback func(c context.Context, rsp interface{}), interceptor entity.ServerInterceptor) {
	ret := func(rsp *PingPong, err error) {
		if err != nil {
			_ = err
		}
		callback(ctx, rsp)
	}
	_GateEntity_Ping_Handler(svc, ctx, in.(*PingPong), ret, interceptor)
	return
}

func _GateEntity_Ping_Remote_Handler(svc interface{}, ctx context.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}
	in := new(PingPong)

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
			MethodName:   "RegisterClient",
			Handler:      _GateEntity_RegisterClient_Remote_Handler,
			LocalHandler: _GateEntity_RegisterClient_Local_Handler,
		},

		"UnregisterClient": {
			MethodName:   "UnregisterClient",
			Handler:      _GateEntity_UnregisterClient_Remote_Handler,
			LocalHandler: _GateEntity_UnregisterClient_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _GateEntity_Ping_Remote_Handler,
			LocalHandler: _GateEntity_Ping_Local_Handler,
		},
	},

	Metadata: "gate.proto",
}
