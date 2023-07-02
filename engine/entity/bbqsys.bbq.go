// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package entity

import (
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"

	// entity "github.com/0x00b/gobbq/proto/entity"

)

var _ = snowflake.GenUUID()

func RegisterBbqSysEntity(etyMgr *EntityManager, impl BbqSysEntity) {
	etyMgr.RegisterEntityDesc(&BbqSysEntityDesc, impl)
}

func NewBbqSysClient(eid EntityID) *BbqSys {
	t := &BbqSys{
		EntityID: eid,
	}
	return t
}

func NewBbqSys(c Context) (*BbqSys, error) {
	etyMgr := GetEntityMgr(c)
	return NewBbqSysWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewBbqSysWithID(c Context, id EntityID) (*BbqSys, error) {

	etyMgr := GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, BbqSysEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil, err
	}
	t := &BbqSys{
		EntityID: id,
	}

	return t, nil
}

type BbqSys struct {
	EntityID EntityID
}

func (t *BbqSys) SysWatch(c Context, req *WatchRequest) (*WatchResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SysWatch"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(WatchResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *BbqSys) SysUnwatch(c Context, req *WatchRequest) (*WatchResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SysUnwatch"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(WatchResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *BbqSys) SysNotify(c Context, req *WatchRequest) (*WatchResponse, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_Unary
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = ""
	pkt.Header.Method = "SysNotify"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.Flags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	defer close(chanRsp)

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return nil, erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, chanRsp)
	if err != nil {
		if !NotMyMethod(err) {
			return nil, err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return nil, err
		}

		pkt.WriteBody(hdrBytes)

		// register callback first, than SendPacket
		RegisterCallback(c, pkt.Header.RequestId, func(pkt *nets.Packet) {
			if pkt.Header.ErrCode != 0 {
				chanRsp <- error(erro.NewError(erro.ErrBadCall.ErrCode, pkt.Header.ErrMsg))
				return
			}
			rsp := new(WatchResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return nil, err
		}
	}

	var rsp any
	select {
	case <-c.Done():
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrContextDone
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		PopCallback(c, pkt.Header.RequestId)
		return nil, erro.ErrTimeOut
	case rsp = <-chanRsp:
	}

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// BbqSysEntity
type BbqSysEntity interface {
	IEntity

	// SysWatch
	SysWatch(c Context, req *WatchRequest) (*WatchResponse, error)

	// SysUnwatch
	SysUnwatch(c Context, req *WatchRequest) (*WatchResponse, error)

	// SysNotify
	SysNotify(c Context, req *WatchRequest) (*WatchResponse, error)
}

func _BbqSysEntity_SysWatch_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysWatch(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysWatch",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysWatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysWatch_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysWatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysWatch_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _BbqSysEntity_SysWatch_Handler(svc, ctx, in, interceptor)
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

func _BbqSysEntity_SysUnwatch_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysUnwatch(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysUnwatch",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysUnwatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysUnwatch_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysUnwatch_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in, interceptor)
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

func _BbqSysEntity_SysNotify_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysNotify(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysNotify",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysNotify(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysNotify_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysNotify_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysNotify_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	npkt := nets.NewPacket()
	defer npkt.Release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.CallType = hdr.CallType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
	npkt.Header.Type = hdr.Type
	npkt.Header.Method = hdr.Method
	npkt.Header.ContentType = hdr.ContentType
	npkt.Header.CompressType = hdr.CompressType
	npkt.Header.Flags = 0
	npkt.Header.TransInfo = hdr.TransInfo

	var rsp any
	if err == nil {
		rsp, err = _BbqSysEntity_SysNotify_Handler(svc, ctx, in, interceptor)
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

var BbqSysEntityDesc = EntityDesc{
	TypeName:    "BbqSysEntity",
	HandlerType: (*BbqSysEntity)(nil),
	Methods: map[string]MethodDesc{

		"SysWatch": {
			MethodName:   "SysWatch",
			Handler:      _BbqSysEntity_SysWatch_Remote_Handler,
			LocalHandler: _BbqSysEntity_SysWatch_Local_Handler,
		},

		"SysUnwatch": {
			MethodName:   "SysUnwatch",
			Handler:      _BbqSysEntity_SysUnwatch_Remote_Handler,
			LocalHandler: _BbqSysEntity_SysUnwatch_Local_Handler,
		},

		"SysNotify": {
			MethodName:   "SysNotify",
			Handler:      _BbqSysEntity_SysNotify_Remote_Handler,
			LocalHandler: _BbqSysEntity_SysNotify_Local_Handler,
		},
	},

	Metadata: "bbqsys.proto",
}
