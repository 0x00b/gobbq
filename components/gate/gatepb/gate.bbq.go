// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package gatepb

import (
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
	// gatepb "github.com/0x00b/gobbq/components/gate/gatepb"
)

var _ = snowflake.GenUUID()

func RegisterGateService(etyMgr *entity.EntityManager, impl GateService) {
	etyMgr.RegisterService(&GateServiceDesc, impl)
}

func NewGateClient() *Gate {
	t := &Gate{}
	return t
}

type Gate struct {
}

func (t *Gate) RegisterClient(c entity.Context, req *RegisterClientRequest) (*RegisterClientResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = GateServiceDesc.TypeName
	pkt.Header.Method = "RegisterClient"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
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

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(RegisterClientResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*RegisterClientResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *Gate) UnregisterClient(c entity.Context, req *RegisterClientRequest) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = GateServiceDesc.TypeName
	pkt.Header.Method = "UnregisterClient"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, nil)
	if err != nil {
		if !entity.NotMyMethod(err) {
			return err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return err
		}

		pkt.WriteBody(hdrBytes)

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *Gate) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = 0
	pkt.Header.Type = GateServiceDesc.TypeName
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
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

		// register callback first, than SendPacket
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// GateService
type GateService interface {
	entity.IService

	// RegisterClient
	RegisterClient(c entity.Context, req *RegisterClientRequest) (*RegisterClientResponse, error)

	// UnregisterClient
	UnregisterClient(c entity.Context, req *RegisterClientRequest) error

	// Ping
	Ping(c entity.Context, req *PingPong) (*PingPong, error)
}

func _GateService_RegisterClient_Handler(svc any, ctx entity.Context, in *RegisterClientRequest, interceptor entity.ServerInterceptor) (*RegisterClientResponse, error) {
	if interceptor == nil {
		return svc.(GateService).RegisterClient(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/RegisterClient",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(GateService).RegisterClient(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterClientResponse), err

}

func _GateService_RegisterClient_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _GateService_RegisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), interceptor)

}

func _GateService_RegisterClient_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterClientRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _GateService_RegisterClient_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

func _GateService_UnregisterClient_Handler(svc any, ctx entity.Context, in *RegisterClientRequest, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(GateService).UnregisterClient(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/UnregisterClient",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(GateService).UnregisterClient(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _GateService_UnregisterClient_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _GateService_UnregisterClient_Handler(svc, ctx, in.(*RegisterClientRequest), interceptor)

}

func _GateService_UnregisterClient_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterClientRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	if err != nil {
		// report
		return
	}
	err = _GateService_UnregisterClient_Handler(svc, ctx, in, interceptor)
	// report err

}

func _GateService_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {
		return svc.(GateService).Ping(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/gatepb.GateService/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(GateService).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*PingPong), err

}

func _GateService_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _GateService_Ping_Handler(svc, ctx, in.(*PingPong), interceptor)

}

func _GateService_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.CheckFlags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _GateService_Ping_Handler(svc, ctx, in, interceptor)
	}
	if err != nil {
		if x, ok := err.(erro.CodeError); ok {
			npkt.Header.ErrCode = x.Code()
			npkt.Header.ErrMsg = x.Message()
		} else {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		}
		npkt.WriteBody(nil)
	} else {
		var rb []byte
		rb, err = codec.DefaultCodec.Marshal(rsp)
		if err != nil {
			npkt.Header.ErrCode = -1
			npkt.Header.ErrMsg = err.Error()
		} else {
			npkt.WriteBody(rb)
		}
	}
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		// report
		_ = err
		return
	}

}

var GateServiceDesc = entity.EntityDesc{
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
