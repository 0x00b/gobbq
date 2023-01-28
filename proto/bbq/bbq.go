// NOTE:!!
//
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package bbq

import (
	"context"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto"
	"github.com/0x00b/gobbq/tool/snowflake"
	// bbq "github.com/0x00b/gobbq/bbq"
)

func RegisterProxyService(impl ProxyService) {
	entity.Manager.RegisterService(&ProxyServiceDesc, impl)
}

func NewProxyService() ProxyService {
	t := &proxyService{}
	return t
}

type proxyService struct {
	entity.Service
}

// RegisterEntity
func (t *proxyService) RegisterEntity(c context.Context, req *RegisterEntityRequest) (rsp *RegisterEntityResponse, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.ProxyService/RegisterEntity",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallService,
		SrcEntity:  nil,
		DstEntity:  nil,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

func (t *proxyService) UnregisterEntity(c context.Context, req *RegisterEntityRequest) (rsp *RegisterEntityResponse, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.ProxyService/UnregisterEntity",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallService,
		SrcEntity:  nil,
		DstEntity:  nil,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

func (t *proxyService) Ping(c context.Context, req *PingPong) (rsp *PingPong, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.ProxyService/Ping",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallService,
		SrcEntity:  nil,
		DstEntity:  nil,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

// ProxyService
type ProxyService interface {
	entity.IService

	// RegisterEntity
	RegisterEntity(c context.Context, req *RegisterEntityRequest) (rsp *RegisterEntityResponse, err error)

	// UnregisterEntity
	UnregisterEntity(c context.Context, req *RegisterEntityRequest) (rsp *RegisterEntityResponse, err error)

	// Ping
	Ping(c context.Context, req *PingPong) (rsp *PingPong, err error)
}

func _ProxyService_RegisterEntity_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(ProxyService).RegisterEntity(ctx, in)
	}
	return nil, nil
}

func _ProxyService_UnregisterEntity_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(ProxyService).UnregisterEntity(ctx, in)
	}
	return nil, nil
}

func _ProxyService_Ping_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingPong)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(ProxyService).Ping(ctx, in)
	}
	return nil, nil
}

var ProxyServiceDesc = entity.ServiceDesc{
	TypeName:    "bbq.ProxyService",
	HandlerType: (*ProxyService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterEntity": {
			MethodName: "RegisterEntity",
			Handler:    _ProxyService_RegisterEntity_Handler,
		},

		"UnregisterEntity": {
			MethodName: "UnregisterEntity",
			Handler:    _ProxyService_UnregisterEntity_Handler,
		},

		"Ping": {
			MethodName: "Ping",
			Handler:    _ProxyService_Ping_Handler,
		},
	},

	Metadata: "bbq.proto",
}

func RegisterGateEntity(impl GateEntity) {
	entity.Manager.RegisterEntity(&GateEntityDesc, impl)
}

func NewGateEntity() GateEntity {
	return NewGateEntityWithID(entity.EntityID(snowflake.GenUUID()))
}

func NewGateEntityWithID(id entity.EntityID) GateEntity {

	ety := entity.NewEntity(id, GateEntityDesc.TypeName)
	t := &gateEntity{entity: ety}

	return t
}

type gateEntity struct {
	entity.EntityClient
	entity *proto.Entity
}

func (t *gateEntity) RegisterClient(c context.Context, req *RegisterClientRequest) (rsp *RegisterClientResponse, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.GateEntity/RegisterClient",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallEntity,
		SrcEntity:  nil,
		DstEntity:  t.entity,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

func (t *gateEntity) UnregisterClient(c context.Context, req *RegisterClientRequest) (rsp *RegisterClientResponse, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.GateEntity/UnregisterClient",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallEntity,
		SrcEntity:  nil,
		DstEntity:  t.entity,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

func (t *gateEntity) Ping(c context.Context, req *PingPong) (rsp *PingPong, err error) {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "bbq.GateEntity/Ping",
		// TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallEntity,
		SrcEntity:  nil,
		DstEntity:  t.entity,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(req)
	if err != nil {
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	//todo get response

	return nil, nil

}

// GateEntity
type GateEntity interface {
	entity.IEntity

	// RegisterClient
	RegisterClient(c context.Context, req *RegisterClientRequest) (rsp *RegisterClientResponse, err error)

	// UnregisterClient
	UnregisterClient(c context.Context, req *RegisterClientRequest) (rsp *RegisterClientResponse, err error)

	// Ping
	Ping(c context.Context, req *PingPong) (rsp *PingPong, err error)
}

func _GateEntity_RegisterClient_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(GateEntity).RegisterClient(ctx, in)
	}
	return nil, nil
}

func _GateEntity_UnregisterClient_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(GateEntity).UnregisterClient(ctx, in)
	}
	return nil, nil
}

func _GateEntity_Ping_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingPong)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(GateEntity).Ping(ctx, in)
	}
	return nil, nil
}

var GateEntityDesc = entity.ServiceDesc{
	TypeName:    "bbq.GateEntity",
	HandlerType: (*GateEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterClient": {
			MethodName: "RegisterClient",
			Handler:    _GateEntity_RegisterClient_Handler,
		},

		"UnregisterClient": {
			MethodName: "UnregisterClient",
			Handler:    _GateEntity_UnregisterClient_Handler,
		},

		"Ping": {
			MethodName: "Ping",
			Handler:    _GateEntity_Ping_Handler,
		},
	},

	Metadata: "bbq.proto",
}
