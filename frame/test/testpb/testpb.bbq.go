// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package testpb

import (
	"errors"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// testpb "github.com/0x00b/gobbq/example/exampb"

)

var _ = snowflake.GenUUID()

func RegisterFrameService(etyMgr *entity.EntityManager, impl FrameService) {
	etyMgr.RegisterService(&FrameServiceDesc, impl)
}

func NewFrameServiceClient() *frameService {
	t := &frameService{}
	return t
}

type frameService struct {
}

func (t *frameService) StartFrame(c entity.Context, req *StartFrameReq) (*StartFrameRsp, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 10
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "testpb.FrameService"}
	pkt.Header.Method = "StartFrame"
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
			rsp := new(StartFrameRsp)
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

	if rsp, ok := rsp.(*StartFrameRsp); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// FrameService
type FrameService interface {
	entity.IService

	// StartFrame
	StartFrame(c entity.Context, req *StartFrameReq) (*StartFrameRsp, error)
}

func _FrameService_StartFrame_Handler(svc any, ctx entity.Context, in *StartFrameReq, interceptor entity.ServerInterceptor) (*StartFrameRsp, error) {
	if interceptor == nil {
		return svc.(FrameService).StartFrame(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/testpb.FrameService/StartFrame",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(FrameService).StartFrame(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*StartFrameRsp), err

}

func _FrameService_StartFrame_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _FrameService_StartFrame_Handler(svc, ctx, in.(*StartFrameReq), interceptor)

}

func _FrameService_StartFrame_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(StartFrameReq)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _FrameService_StartFrame_Handler(svc, ctx, in, interceptor)

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

var FrameServiceDesc = entity.EntityDesc{
	TypeName:    "testpb.FrameService",
	HandlerType: (*FrameService)(nil),
	Methods: map[string]entity.MethodDesc{

		"StartFrame": {
			MethodName:   "StartFrame",
			Handler:      _FrameService_StartFrame_Remote_Handler,
			LocalHandler: _FrameService_StartFrame_Local_Handler,
		},
	},

	Metadata: "testpb.proto",
}
