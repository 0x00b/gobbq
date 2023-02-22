// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

package proxypb

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"

	// proxypb "github.com/0x00b/gobbq/components/proxy/proxypb"

)

var _ = snowflake.GenUUID()

func RegisterProxyService(impl ProxyService) {
	entity.Manager.RegisterService(&ProxyServiceDesc, impl)
}

func NewProxyServiceClient(client *codec.PacketReadWriter) *proxyService {
	t := &proxyService{client: client}
	return t
}

func NewProxyService(client *codec.PacketReadWriter) *proxyService {
	t := &proxyService{client: client}
	return t
}

type proxyService struct {
	client *codec.PacketReadWriter
}

func (t *proxyService) RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "RegisterProxy"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterProxyResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterProxyResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "RegisterInst"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterInstResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterInstResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "SyncService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SyncServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*SyncServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) RegisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "RegisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterEntityResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterEntityResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) RegisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "RegisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) UnregisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "UnregisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterEntityResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterEntityResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) UnregisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "UnregisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyService) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxyService"}
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	// err = entity.HandleCallLocalMethod(pkt, req, itfCallback)
	// if err == nil {
	// 	return nil
	// }

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(req)
	if err != nil {
		xlog.Errorln(err)
		return nil, err
	}

	pkt.WriteBody(hdrBytes)

	t.client.WritePacket(pkt)

	// register callback
	chanRsp := make(chan any)
	if c != nil {
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
			close(chanRsp)
		})
	}
	rsp := <-chanRsp
	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ProxyService
type ProxyService interface {
	entity.IService

	// RegisterProxy
	RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error)

	// RegisterInst
	RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error)

	// SyncService
	SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error)

	// RegisterEntity
	RegisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error)

	// RegisterService
	RegisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error)

	// UnregisterEntity
	UnregisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error)

	// UnregisterService
	UnregisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error)

	// Ping
	Ping(c entity.Context, req *PingPong) (*PingPong, error)
}

func _ProxyService_RegisterProxy_Handler(svc any, ctx entity.Context, in *RegisterProxyRequest, interceptor entity.ServerInterceptor) (*RegisterProxyResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).RegisterProxy(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterProxy",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).RegisterProxy(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterProxyResponse), err

}

//func _ProxyService_RegisterProxy_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterProxyResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_RegisterProxy_Handler(svc, ctx, in.(*RegisterProxyRequest) , ret, interceptor)
//
//}

func _ProxyService_RegisterProxy_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterProxyRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_RegisterProxy_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_RegisterInst_Handler(svc any, ctx entity.Context, in *RegisterInstRequest, interceptor entity.ServerInterceptor) (*RegisterInstResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).RegisterInst(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterInst",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).RegisterInst(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterInstResponse), err

}

//func _ProxyService_RegisterInst_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterInstResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_RegisterInst_Handler(svc, ctx, in.(*RegisterInstRequest) , ret, interceptor)
//
//}

func _ProxyService_RegisterInst_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterInstRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_RegisterInst_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_SyncService_Handler(svc any, ctx entity.Context, in *SyncServiceRequest, interceptor entity.ServerInterceptor) (*SyncServiceResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).SyncService(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/SyncService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).SyncService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*SyncServiceResponse), err

}

//func _ProxyService_SyncService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *SyncServiceResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_SyncService_Handler(svc, ctx, in.(*SyncServiceRequest) , ret, interceptor)
//
//}

func _ProxyService_SyncService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SyncServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_SyncService_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_RegisterEntity_Handler(svc any, ctx entity.Context, in *RegisterEntityRequest, interceptor entity.ServerInterceptor) (*RegisterEntityResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).RegisterEntity(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterEntity",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).RegisterEntity(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterEntityResponse), err

}

//func _ProxyService_RegisterEntity_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterEntityResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_RegisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest) , ret, interceptor)
//
//}

func _ProxyService_RegisterEntity_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_RegisterEntity_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_RegisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).RegisterService(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/RegisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).RegisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterServiceResponse), err

}

//func _ProxyService_RegisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterServiceResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_RegisterService_Handler(svc, ctx, in.(*RegisterServiceRequest) , ret, interceptor)
//
//}

func _ProxyService_RegisterService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_RegisterService_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_UnregisterEntity_Handler(svc any, ctx entity.Context, in *RegisterEntityRequest, interceptor entity.ServerInterceptor) (*RegisterEntityResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).UnregisterEntity(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/UnregisterEntity",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).UnregisterEntity(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterEntityResponse), err

}

//func _ProxyService_UnregisterEntity_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterEntityResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_UnregisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest) , ret, interceptor)
//
//}

func _ProxyService_UnregisterEntity_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_UnregisterEntity_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_UnregisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {

		return svc.(ProxyService).UnregisterService(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/UnregisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).UnregisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*RegisterServiceResponse), err

}

//func _ProxyService_UnregisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *RegisterServiceResponse, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_UnregisterService_Handler(svc, ctx, in.(*RegisterServiceRequest) , ret, interceptor)
//
//}

func _ProxyService_UnregisterService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_UnregisterService_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

func _ProxyService_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {

		return svc.(ProxyService).Ping(ctx, in)

	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyService/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyService).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	return rsp.(*PingPong), err

}

//func _ProxyService_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor)(any, error) {
//
//		ret := func(rsp *PingPong, err error) {
//			if err != nil {
//				_ = err
//			}
//			callback(ctx, rsp)
//		}
//
//
//	_ProxyService_Ping_Handler(svc, ctx, in.(*PingPong) , ret, interceptor)
//
//}

func _ProxyService_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyService_Ping_Handler(svc, ctx, in, interceptor)

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
	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		xlog.Errorln("WritePacket", err)
		return
	}

}

var ProxyServiceDesc = entity.EntityDesc{
	TypeName:    "proxypb.ProxyService",
	HandlerType: (*ProxyService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterProxy": {
			MethodName: "RegisterProxy",
			Handler:    _ProxyService_RegisterProxy_Remote_Handler,
			//LocalHandler:	_ProxyService_RegisterProxy_Local_Handler,
		},

		"RegisterInst": {
			MethodName: "RegisterInst",
			Handler:    _ProxyService_RegisterInst_Remote_Handler,
			//LocalHandler:	_ProxyService_RegisterInst_Local_Handler,
		},

		"SyncService": {
			MethodName: "SyncService",
			Handler:    _ProxyService_SyncService_Remote_Handler,
			//LocalHandler:	_ProxyService_SyncService_Local_Handler,
		},

		"RegisterEntity": {
			MethodName: "RegisterEntity",
			Handler:    _ProxyService_RegisterEntity_Remote_Handler,
			//LocalHandler:	_ProxyService_RegisterEntity_Local_Handler,
		},

		"RegisterService": {
			MethodName: "RegisterService",
			Handler:    _ProxyService_RegisterService_Remote_Handler,
			//LocalHandler:	_ProxyService_RegisterService_Local_Handler,
		},

		"UnregisterEntity": {
			MethodName: "UnregisterEntity",
			Handler:    _ProxyService_UnregisterEntity_Remote_Handler,
			//LocalHandler:	_ProxyService_UnregisterEntity_Local_Handler,
		},

		"UnregisterService": {
			MethodName: "UnregisterService",
			Handler:    _ProxyService_UnregisterService_Remote_Handler,
			//LocalHandler:	_ProxyService_UnregisterService_Local_Handler,
		},

		"Ping": {
			MethodName: "Ping",
			Handler:    _ProxyService_Ping_Remote_Handler,
			//LocalHandler:	_ProxyService_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}
