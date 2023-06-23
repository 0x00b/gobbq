// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package bbqsys

import (
	"errors"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// bbqsys "github.com/0x00b/gobbq/proto/bbqsys"

)

var _ = snowflake.GenUUID()

func RegisterBbqSysEntity(etyMgr *entity.EntityManager, impl BbqSysEntity) {
	etyMgr.RegisterEntityDesc(&BbqSysEntityDesc, impl)
}

func NewBbqSysEntityClient(eid entity.EntityID) *bbqSysEntity {
	t := &bbqSysEntity{
		EntityID: eid,
	}
	return t
}

func NewBbqSysEntity(c entity.Context) *bbqSysEntity {
	etyMgr := entity.GetEntityMgr(c)
	return NewBbqSysEntityWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewBbqSysEntityWithID(c entity.Context, id entity.EntityID) *bbqSysEntity {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, BbqSysEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &bbqSysEntity{
		EntityID: id,
	}

	return t
}

type bbqSysEntity struct {
	EntityID entity.EntityID
}

func (t *bbqSysEntity) SysWatch(c entity.Context, req *WatchRequest) (*WatchResponse, error) {

	pkt, release := nets.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = BbqSysEntityDesc.TypeName
	pkt.Header.Method = "SysWatch"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
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
			rsp := new(WatchResponse)
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
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	close(chanRsp)

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *bbqSysEntity) SysUnwatch(c entity.Context, req *WatchRequest) (*WatchResponse, error) {

	pkt, release := nets.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = BbqSysEntityDesc.TypeName
	pkt.Header.Method = "SysUnwatch"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
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
			rsp := new(WatchResponse)
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
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	close(chanRsp)

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *bbqSysEntity) SysNotify(c entity.Context, req *WatchRequest) (*WatchResponse, error) {

	pkt, release := nets.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = BbqSysEntityDesc.TypeName
	pkt.Header.Method = "SysNotify"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)
	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return nil, errors.New("bad context")
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
			rsp := new(WatchResponse)
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
		return nil, errors.New("context done")
	case <-time.After(time.Duration(pkt.Header.Timeout) * time.Second):
		entity.PopCallback(c, pkt.Header.RequestId)
		return nil, errors.New("time out")
	case rsp = <-chanRsp:
	}

	close(chanRsp)

	if rsp, ok := rsp.(*WatchResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// BbqSysEntity
type BbqSysEntity interface {
	entity.IEntity

	// SysWatch
	SysWatch(c entity.Context, req *WatchRequest) (*WatchResponse, error)

	// SysUnwatch
	SysUnwatch(c entity.Context, req *WatchRequest) (*WatchResponse, error)

	// SysNotify
	SysNotify(c entity.Context, req *WatchRequest) (*WatchResponse, error)
}

func _BbqSysEntity_SysWatch_Handler(svc any, ctx entity.Context, in *WatchRequest, interceptor entity.ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysWatch(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/bbqsys.BbqSysEntity/SysWatch",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysWatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysWatch_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysWatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysWatch_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _BbqSysEntity_SysWatch_Handler(svc, ctx, in, interceptor)

	npkt, release := nets.NewPacket()
	defer release()

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
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _BbqSysEntity_SysUnwatch_Handler(svc any, ctx entity.Context, in *WatchRequest, interceptor entity.ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysUnwatch(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/bbqsys.BbqSysEntity/SysUnwatch",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysUnwatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysUnwatch_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysUnwatch_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in, interceptor)

	npkt, release := nets.NewPacket()
	defer release()

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
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

func _BbqSysEntity_SysNotify_Handler(svc any, ctx entity.Context, in *WatchRequest, interceptor entity.ServerInterceptor) (*WatchResponse, error) {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysNotify(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/bbqsys.BbqSysEntity/SysNotify",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(BbqSysEntity).SysNotify(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*WatchResponse), err

}

func _BbqSysEntity_SysNotify_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _BbqSysEntity_SysNotify_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysNotify_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _BbqSysEntity_SysNotify_Handler(svc, ctx, in, interceptor)

	npkt, release := nets.NewPacket()
	defer release()

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
	err = pkt.Src.SendPacket(npkt)
	if err != nil {
		xlog.Errorln("SendPacket", err)
		return
	}

}

var BbqSysEntityDesc = entity.EntityDesc{
	TypeName:    "bbqsys.BbqSysEntity",
	HandlerType: (*BbqSysEntity)(nil),
	Methods: map[string]entity.MethodDesc{

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

	Metadata: "sys.proto",
}
