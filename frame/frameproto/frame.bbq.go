// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package frameproto

import (
	"errors"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// frameproto "github.com/0x00b/gobbq/frame/frameproto"

)

var _ = snowflake.GenUUID()

func RegisterFrameSeverEntity(etyMgr *entity.EntityManager, impl FrameSeverEntity) {
	etyMgr.RegisterEntityDesc(&FrameSeverEntityDesc, impl)
}

func NewFrameSeverClient(eid entity.EntityID) *FrameSever {
	t := &FrameSever{
		EntityID: eid,
	}
	return t
}

func NewFrameSever(c entity.Context) *FrameSever {
	etyMgr := entity.GetEntityMgr(c)
	return NewFrameSeverWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewFrameSeverWithID(c entity.Context, id entity.EntityID) *FrameSever {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, FrameSeverEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &FrameSever{
		EntityID: id,
	}

	return t
}

type FrameSever struct {
	EntityID entity.EntityID
}

func (t *FrameSever) Heartbeat(c entity.Context, req *HeartbeatReq) (*HeartbeatRsp, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Heartbeat"
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
			rsp := new(HeartbeatRsp)
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

	if rsp, ok := rsp.(*HeartbeatRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *FrameSever) Init(c entity.Context, req *InitReq) (*InitRsp, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Init"
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
			rsp := new(InitRsp)
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

	if rsp, ok := rsp.(*InitRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *FrameSever) Join(c entity.Context, req *JoinReq) (*JoinRsp, error) {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Join"
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
			rsp := new(JoinRsp)
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

	if rsp, ok := rsp.(*JoinRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *FrameSever) Progress(c entity.Context, req *ProgressReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Progress"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameSever) Ready(c entity.Context, req *ReadyReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Ready"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameSever) Input(c entity.Context, req *InputReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "Input"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameSever) GameOver(c entity.Context, req *GameOverReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameSeverEntityDesc.TypeName
	pkt.Header.Method = "GameOver"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

// FrameSeverEntity
type FrameSeverEntity interface {
	entity.IEntity

	// Heartbeat
	Heartbeat(c entity.Context, req *HeartbeatReq) (*HeartbeatRsp, error)

	// Init
	Init(c entity.Context, req *InitReq) (*InitRsp, error)

	// Join 客户端加入
	Join(c entity.Context, req *JoinReq) (*JoinRsp, error)

	// Progress 客户端上报加载进度
	Progress(c entity.Context, req *ProgressReq) error

	// Ready 客户端准备好
	Ready(c entity.Context, req *ReadyReq) error

	// Input 客户端的操作
	Input(c entity.Context, req *InputReq) error

	// GameOver 上报结束
	GameOver(c entity.Context, req *GameOverReq) error
}

func _FrameSeverEntity_Heartbeat_Handler(svc any, ctx entity.Context, in *HeartbeatReq, interceptor entity.ServerInterceptor) (*HeartbeatRsp, error) {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Heartbeat(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Heartbeat",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(FrameSeverEntity).Heartbeat(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*HeartbeatRsp), err

}

func _FrameSeverEntity_Heartbeat_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _FrameSeverEntity_Heartbeat_Handler(svc, ctx, in.(*HeartbeatReq), interceptor)

}

func _FrameSeverEntity_Heartbeat_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(HeartbeatReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Heartbeat_Handler(svc, ctx, in, interceptor)

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

func _FrameSeverEntity_Init_Handler(svc any, ctx entity.Context, in *InitReq, interceptor entity.ServerInterceptor) (*InitRsp, error) {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Init(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Init",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(FrameSeverEntity).Init(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*InitRsp), err

}

func _FrameSeverEntity_Init_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _FrameSeverEntity_Init_Handler(svc, ctx, in.(*InitReq), interceptor)

}

func _FrameSeverEntity_Init_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(InitReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Init_Handler(svc, ctx, in, interceptor)

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

func _FrameSeverEntity_Join_Handler(svc any, ctx entity.Context, in *JoinReq, interceptor entity.ServerInterceptor) (*JoinRsp, error) {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Join(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Join",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(FrameSeverEntity).Join(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*JoinRsp), err

}

func _FrameSeverEntity_Join_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _FrameSeverEntity_Join_Handler(svc, ctx, in.(*JoinReq), interceptor)

}

func _FrameSeverEntity_Join_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(JoinReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Join_Handler(svc, ctx, in, interceptor)

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

func _FrameSeverEntity_Progress_Handler(svc any, ctx entity.Context, in *ProgressReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Progress(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Progress",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Progress(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Progress_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Progress_Handler(svc, ctx, in.(*ProgressReq), interceptor)

}

func _FrameSeverEntity_Progress_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ProgressReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Progress_Handler(svc, ctx, in, interceptor)

}

func _FrameSeverEntity_Ready_Handler(svc any, ctx entity.Context, in *ReadyReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Ready(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Ready",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Ready(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Ready_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Ready_Handler(svc, ctx, in.(*ReadyReq), interceptor)

}

func _FrameSeverEntity_Ready_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ReadyReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Ready_Handler(svc, ctx, in, interceptor)

}

func _FrameSeverEntity_Input_Handler(svc any, ctx entity.Context, in *InputReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Input(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Input",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Input(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Input_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Input_Handler(svc, ctx, in.(*InputReq), interceptor)

}

func _FrameSeverEntity_Input_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(InputReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Input_Handler(svc, ctx, in, interceptor)

}

func _FrameSeverEntity_GameOver_Handler(svc any, ctx entity.Context, in *GameOverReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).GameOver(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/GameOver",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).GameOver(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_GameOver_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_GameOver_Handler(svc, ctx, in.(*GameOverReq), interceptor)

}

func _FrameSeverEntity_GameOver_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(GameOverReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_GameOver_Handler(svc, ctx, in, interceptor)

}

var FrameSeverEntityDesc = entity.EntityDesc{
	TypeName:    "frameproto.FrameSeverEntity",
	HandlerType: (*FrameSeverEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"Heartbeat": {
			MethodName:   "Heartbeat",
			Handler:      _FrameSeverEntity_Heartbeat_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Heartbeat_Local_Handler,
		},

		"Init": {
			MethodName:   "Init",
			Handler:      _FrameSeverEntity_Init_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Init_Local_Handler,
		},

		"Join": {
			MethodName:   "Join",
			Handler:      _FrameSeverEntity_Join_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Join_Local_Handler,
		},

		"Progress": {
			MethodName:   "Progress",
			Handler:      _FrameSeverEntity_Progress_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Progress_Local_Handler,
		},

		"Ready": {
			MethodName:   "Ready",
			Handler:      _FrameSeverEntity_Ready_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Ready_Local_Handler,
		},

		"Input": {
			MethodName:   "Input",
			Handler:      _FrameSeverEntity_Input_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Input_Local_Handler,
		},

		"GameOver": {
			MethodName:   "GameOver",
			Handler:      _FrameSeverEntity_GameOver_Remote_Handler,
			LocalHandler: _FrameSeverEntity_GameOver_Local_Handler,
		},
	},

	Metadata: "frame.proto",
}

func RegisterFrameClientEntity(etyMgr *entity.EntityManager, impl FrameClientEntity) {
	etyMgr.RegisterEntityDesc(&FrameClientEntityDesc, impl)
}

func NewFrameClientClient(eid entity.EntityID) *FrameClient {
	t := &FrameClient{
		EntityID: eid,
	}
	return t
}

func NewFrameClient(c entity.Context) *FrameClient {
	etyMgr := entity.GetEntityMgr(c)
	return NewFrameClientWithID(c, etyMgr.EntityIDGenerator.NewEntityID())
}

func NewFrameClientWithID(c entity.Context, id entity.EntityID) *FrameClient {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id, FrameClientEntityDesc.TypeName)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &FrameClient{
		EntityID: id,
	}

	return t
}

type FrameClient struct {
	EntityID entity.EntityID
}

func (t *FrameClient) Progress(c entity.Context, req *ProgressReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameClientEntityDesc.TypeName
	pkt.Header.Method = "Progress"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameClient) Start(c entity.Context, req *StartReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameClientEntityDesc.TypeName
	pkt.Header.Method = "Start"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameClient) Frame(c entity.Context, req *FrameReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameClientEntityDesc.TypeName
	pkt.Header.Method = "Frame"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

func (t *FrameClient) GameOver(c entity.Context, req *GameOverReq) error {

	pkt := nets.NewPacket()
	defer pkt.Release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = uint64(c.EntityID())
	pkt.Header.DstEntity = uint64(t.EntityID)
	pkt.Header.Type = FrameClientEntityDesc.TypeName
	pkt.Header.Method = "GameOver"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	etyMgr := entity.GetEntityMgr(c)
	if etyMgr == nil {
		return errors.New("bad context")
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

// FrameClientEntity
type FrameClientEntity interface {
	entity.IEntity

	// Progress 通知客户端其他人加载进度
	Progress(c entity.Context, req *ProgressReq) error

	// Start 通知客户端开始
	Start(c entity.Context, req *StartReq) error

	// Frame 帧下发
	Frame(c entity.Context, req *FrameReq) error

	// GameOver 游戏结束
	GameOver(c entity.Context, req *GameOverReq) error
}

func _FrameClientEntity_Progress_Handler(svc any, ctx entity.Context, in *ProgressReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).Progress(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/Progress",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).Progress(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_Progress_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_Progress_Handler(svc, ctx, in.(*ProgressReq), interceptor)

}

func _FrameClientEntity_Progress_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ProgressReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_Progress_Handler(svc, ctx, in, interceptor)

}

func _FrameClientEntity_Start_Handler(svc any, ctx entity.Context, in *StartReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).Start(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/Start",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).Start(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_Start_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_Start_Handler(svc, ctx, in.(*StartReq), interceptor)

}

func _FrameClientEntity_Start_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(StartReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_Start_Handler(svc, ctx, in, interceptor)

}

func _FrameClientEntity_Frame_Handler(svc any, ctx entity.Context, in *FrameReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).Frame(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/Frame",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).Frame(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_Frame_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_Frame_Handler(svc, ctx, in.(*FrameReq), interceptor)

}

func _FrameClientEntity_Frame_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(FrameReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_Frame_Handler(svc, ctx, in, interceptor)

}

func _FrameClientEntity_GameOver_Handler(svc any, ctx entity.Context, in *GameOverReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).GameOver(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/GameOver",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).GameOver(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_GameOver_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_GameOver_Handler(svc, ctx, in.(*GameOverReq), interceptor)

}

func _FrameClientEntity_GameOver_Remote_Handler(svc any, ctx entity.Context, pkt *nets.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(GameOverReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_GameOver_Handler(svc, ctx, in, interceptor)

}

var FrameClientEntityDesc = entity.EntityDesc{
	TypeName:    "frameproto.FrameClientEntity",
	HandlerType: (*FrameClientEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"Progress": {
			MethodName:   "Progress",
			Handler:      _FrameClientEntity_Progress_Remote_Handler,
			LocalHandler: _FrameClientEntity_Progress_Local_Handler,
		},

		"Start": {
			MethodName:   "Start",
			Handler:      _FrameClientEntity_Start_Remote_Handler,
			LocalHandler: _FrameClientEntity_Start_Local_Handler,
		},

		"Frame": {
			MethodName:   "Frame",
			Handler:      _FrameClientEntity_Frame_Remote_Handler,
			LocalHandler: _FrameClientEntity_Frame_Local_Handler,
		},

		"GameOver": {
			MethodName:   "GameOver",
			Handler:      _FrameClientEntity_GameOver_Remote_Handler,
			LocalHandler: _FrameClientEntity_GameOver_Local_Handler,
		},
	},

	Metadata: "frame.proto",
}
