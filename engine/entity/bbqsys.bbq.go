// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package entity

import (
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

func (t *BbqSys) SysWatch(c Context, req *WatchRequest) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_OneWay
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

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, nil)
	if err != nil {
		if !NotMyMethod(err) {
			return err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return err
		}

		pkt.WriteBody(hdrBytes)

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *BbqSys) SysUnwatch(c Context, req *WatchRequest) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_OneWay
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

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, nil)
	if err != nil {
		if !NotMyMethod(err) {
			return err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return err
		}

		pkt.WriteBody(hdrBytes)

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *BbqSys) SysNotify(c Context, req *WatchRequest) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.CallType = bbq.CallType_OneWay
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

	etyMgr := GetEntityMgr(c)
	if etyMgr == nil {
		return erro.ErrBadContext
	}
	err := etyMgr.LocalCall(pkt, req, nil)
	if err != nil {
		if !NotMyMethod(err) {
			return err
		}

		hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
		if err != nil {
			xlog.Errorln(err)
			return err
		}

		pkt.WriteBody(hdrBytes)

		err = GetProxy(c).SendPacket(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

// BbqSysEntity
type BbqSysEntity interface {
	IEntity

	// SysWatch
	SysWatch(c Context, req *WatchRequest) error

	// SysUnwatch
	SysUnwatch(c Context, req *WatchRequest) error

	// SysNotify
	SysNotify(c Context, req *WatchRequest) error
}

func _BbqSysEntity_SysWatch_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) error {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysWatch(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysWatch",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return nil, svc.(BbqSysEntity).SysWatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _BbqSysEntity_SysWatch_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return nil, _BbqSysEntity_SysWatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysWatch_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	if err != nil {
		// report
		return
	}
	err = _BbqSysEntity_SysWatch_Handler(svc, ctx, in, interceptor)
	_ = err
	// report err

}

func _BbqSysEntity_SysUnwatch_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) error {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysUnwatch(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysUnwatch",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return nil, svc.(BbqSysEntity).SysUnwatch(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _BbqSysEntity_SysUnwatch_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return nil, _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysUnwatch_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	if err != nil {
		// report
		return
	}
	err = _BbqSysEntity_SysUnwatch_Handler(svc, ctx, in, interceptor)
	_ = err
	// report err

}

func _BbqSysEntity_SysNotify_Handler(svc any, ctx Context, in *WatchRequest, interceptor ServerInterceptor) error {
	if interceptor == nil {
		return svc.(BbqSysEntity).SysNotify(ctx, in)
	}

	info := &ServerInfo{
		Server:     svc,
		FullMethod: "/BbqSysEntity/SysNotify",
	}

	handler := func(ctx Context, rsp any) (any, error) {

		return nil, svc.(BbqSysEntity).SysNotify(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _BbqSysEntity_SysNotify_Local_Handler(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error) {

	return nil, _BbqSysEntity_SysNotify_Handler(svc, ctx, in.(*WatchRequest), interceptor)

}

func _BbqSysEntity_SysNotify_Remote_Handler(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor) {

	hdr := pkt.Header

	in := new(WatchRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)

	if err != nil {
		// report
		return
	}
	err = _BbqSysEntity_SysNotify_Handler(svc, ctx, in, interceptor)
	_ = err
	// report err

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
