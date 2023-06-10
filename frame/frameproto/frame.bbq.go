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
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// frameproto "github.com/0x00b/gobbq/frame/frameproto"

)

var _ = snowflake.GenUUID()

func RegisterFrameSeverEntity(etyMgr *entity.EntityManager, impl FrameSeverEntity) {
	etyMgr.RegisterEntityDesc(&FrameSeverEntityDesc, impl)
}

func NewFrameSeverEntityClient(eid *bbq.EntityID) *frameSeverEntity {
	t := &frameSeverEntity{
		EntityID: eid,
	}
	return t
}

func NewFrameSeverEntity(c entity.Context) *frameSeverEntity {
	etyMgr := entity.GetEntityMgr(c)
	return NewFrameSeverEntityWithID(c, etyMgr.EntityIDGenerator.NewEntityID("frameproto.FrameSeverEntity"))
}

func NewFrameSeverEntityWithID(c entity.Context, id *bbq.EntityID) *frameSeverEntity {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &frameSeverEntity{
		EntityID: id,
	}

	return t
}

type frameSeverEntity struct {
	EntityID *bbq.EntityID
}

func (t *frameSeverEntity) Heartbeat(c entity.Context, req *HeartbeatReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Heartbeat"
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameSeverEntity) Init(c entity.Context, req *InitReq) (*InitRsp, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Init"
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

		// register callback first, than SendPackt
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(InitRsp)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
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

	if rsp, ok := rsp.(*InitRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *frameSeverEntity) Join(c entity.Context, req *JoinReq) (*JoinRsp, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Join"
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

		// register callback first, than SendPackt
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(JoinRsp)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
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

	if rsp, ok := rsp.(*JoinRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *frameSeverEntity) Progress(c entity.Context, req *ProgressReq) (*ProgressRsp, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Progress"
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

		// register callback first, than SendPackt
		entity.RegisterCallback(c, pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(ProgressRsp)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
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

	if rsp, ok := rsp.(*ProgressRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *frameSeverEntity) Ready(c entity.Context, req *ReadyReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameSeverEntity) Move(c entity.Context, req *MoveReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Move"
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameSeverEntity) Input(c entity.Context, req *InputReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameSeverEntity) Result(c entity.Context, req *ResultReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Result"
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
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
	Heartbeat(c entity.Context, req *HeartbeatReq) error

	// Init
	Init(c entity.Context, req *InitReq) (*InitRsp, error)

	// Join
	Join(c entity.Context, req *JoinReq) (*JoinRsp, error)

	// Progress
	Progress(c entity.Context, req *ProgressReq) (*ProgressRsp, error)

	// Ready
	Ready(c entity.Context, req *ReadyReq) error

	// Move
	Move(c entity.Context, req *MoveReq) error

	// Input
	Input(c entity.Context, req *InputReq) error

	// Result
	Result(c entity.Context, req *ResultReq) error
}

func _FrameSeverEntity_Heartbeat_Handler(svc any, ctx entity.Context, in *HeartbeatReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Heartbeat(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Heartbeat",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Heartbeat(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Heartbeat_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Heartbeat_Handler(svc, ctx, in.(*HeartbeatReq), interceptor)

}

func _FrameSeverEntity_Heartbeat_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(HeartbeatReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Heartbeat_Handler(svc, ctx, in, interceptor)

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

func _FrameSeverEntity_Init_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(InitReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Init_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.SendPackt(npkt)
	if err != nil {
		xlog.Errorln("SendPackt", err)
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

func _FrameSeverEntity_Join_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(JoinReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Join_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.SendPackt(npkt)
	if err != nil {
		xlog.Errorln("SendPackt", err)
		return
	}

}

func _FrameSeverEntity_Progress_Handler(svc any, ctx entity.Context, in *ProgressReq, interceptor entity.ServerInterceptor) (*ProgressRsp, error) {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Progress(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Progress",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(FrameSeverEntity).Progress(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*ProgressRsp), err

}

func _FrameSeverEntity_Progress_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _FrameSeverEntity_Progress_Handler(svc, ctx, in.(*ProgressReq), interceptor)

}

func _FrameSeverEntity_Progress_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ProgressReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameSeverEntity_Progress_Handler(svc, ctx, in, interceptor)

	npkt, release := codec.NewPacket()
	defer release()

	npkt.Header.Version = hdr.Version
	npkt.Header.RequestId = hdr.RequestId
	npkt.Header.Timeout = hdr.Timeout
	npkt.Header.RequestType = bbq.RequestType_RequestRespone
	npkt.Header.ServiceType = hdr.ServiceType
	npkt.Header.SrcEntity = hdr.DstEntity
	npkt.Header.DstEntity = hdr.SrcEntity
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
	err = pkt.Src.SendPackt(npkt)
	if err != nil {
		xlog.Errorln("SendPackt", err)
		return
	}

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

func _FrameSeverEntity_Ready_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

func _FrameSeverEntity_Move_Handler(svc any, ctx entity.Context, in *MoveReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Move(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Move",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Move(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Move_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Move_Handler(svc, ctx, in.(*MoveReq), interceptor)

}

func _FrameSeverEntity_Move_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(MoveReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Move_Handler(svc, ctx, in, interceptor)

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

func _FrameSeverEntity_Input_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

func _FrameSeverEntity_Result_Handler(svc any, ctx entity.Context, in *ResultReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameSeverEntity).Result(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameSeverEntity/Result",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameSeverEntity).Result(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameSeverEntity_Result_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameSeverEntity_Result_Handler(svc, ctx, in.(*ResultReq), interceptor)

}

func _FrameSeverEntity_Result_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ResultReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameSeverEntity_Result_Handler(svc, ctx, in, interceptor)

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

		"Move": {
			MethodName:   "Move",
			Handler:      _FrameSeverEntity_Move_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Move_Local_Handler,
		},

		"Input": {
			MethodName:   "Input",
			Handler:      _FrameSeverEntity_Input_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Input_Local_Handler,
		},

		"Result": {
			MethodName:   "Result",
			Handler:      _FrameSeverEntity_Result_Remote_Handler,
			LocalHandler: _FrameSeverEntity_Result_Local_Handler,
		},
	},

	Metadata: "frame.proto",
}

func RegisterFrameClientEntity(etyMgr *entity.EntityManager, impl FrameClientEntity) {
	etyMgr.RegisterEntityDesc(&FrameClientEntityDesc, impl)
}

func NewFrameClientEntityClient(eid *bbq.EntityID) *frameClientEntity {
	t := &frameClientEntity{
		EntityID: eid,
	}
	return t
}

func NewFrameClientEntity(c entity.Context) *frameClientEntity {
	etyMgr := entity.GetEntityMgr(c)
	return NewFrameClientEntityWithID(c, etyMgr.EntityIDGenerator.NewEntityID("frameproto.FrameClientEntity"))
}

func NewFrameClientEntityWithID(c entity.Context, id *bbq.EntityID) *frameClientEntity {

	etyMgr := entity.GetEntityMgr(c)
	_, err := etyMgr.NewEntity(c, id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &frameClientEntity{
		EntityID: id,
	}

	return t
}

type frameClientEntity struct {
	EntityID *bbq.EntityID
}

func (t *frameClientEntity) Start(c entity.Context, req *StartReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameClientEntity) Frame(c entity.Context, req *FrameReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameClientEntity) Result(c entity.Context, req *ClientResultReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Result"
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

func (t *frameClientEntity) Close(c entity.Context, req *CloseReq) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.EntityID
	pkt.Header.Method = "Close"
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

		err = entity.GetRemoteEntityManager(c).SendPackt(pkt)
		if err != nil {
			return err
		}
	}

	return nil

}

// FrameClientEntity
type FrameClientEntity interface {
	entity.IEntity

	// Start
	Start(c entity.Context, req *StartReq) error

	// Frame
	Frame(c entity.Context, req *FrameReq) error

	// Result
	Result(c entity.Context, req *ClientResultReq) error

	// Close
	Close(c entity.Context, req *CloseReq) error
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

func _FrameClientEntity_Start_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

func _FrameClientEntity_Frame_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

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

func _FrameClientEntity_Result_Handler(svc any, ctx entity.Context, in *ClientResultReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).Result(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/Result",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).Result(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_Result_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_Result_Handler(svc, ctx, in.(*ClientResultReq), interceptor)

}

func _FrameClientEntity_Result_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(ClientResultReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_Result_Handler(svc, ctx, in, interceptor)

}

func _FrameClientEntity_Close_Handler(svc any, ctx entity.Context, in *CloseReq, interceptor entity.ServerInterceptor) error {
	if interceptor == nil {
		return svc.(FrameClientEntity).Close(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/frameproto.FrameClientEntity/Close",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return nil, svc.(FrameClientEntity).Close(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return err

}

func _FrameClientEntity_Close_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return nil, _FrameClientEntity_Close_Handler(svc, ctx, in.(*CloseReq), interceptor)

}

func _FrameClientEntity_Close_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(CloseReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// err
		return
	}

	_FrameClientEntity_Close_Handler(svc, ctx, in, interceptor)

}

var FrameClientEntityDesc = entity.EntityDesc{
	TypeName:    "frameproto.FrameClientEntity",
	HandlerType: (*FrameClientEntity)(nil),
	Methods: map[string]entity.MethodDesc{

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

		"Result": {
			MethodName:   "Result",
			Handler:      _FrameClientEntity_Result_Remote_Handler,
			LocalHandler: _FrameClientEntity_Result_Local_Handler,
		},

		"Close": {
			MethodName:   "Close",
			Handler:      _FrameClientEntity_Close_Remote_Handler,
			LocalHandler: _FrameClientEntity_Close_Local_Handler,
		},
	},

	Metadata: "frame.proto",
}
