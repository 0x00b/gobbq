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

func RegisterProxyEtyEntity(impl ProxyEtyEntity) {
	entity.Manager.RegisterEntity(&ProxyEtyEntityDesc, impl)
}

func NewProxyEtyEntityClient(client *codec.PacketReadWriter, entity *bbq.EntityID) *proxyEtyEntity {
	t := &proxyEtyEntity{client: client, entity: entity}
	return t
}

func NewProxyEtyEntity(c entity.Context, client *codec.PacketReadWriter) *proxyEtyEntity {
	return NewProxyEtyEntityWithID(c, entity.NewEntityID.NewEntityID("proxypb.ProxyEtyEntity"), client)
}

func NewProxyEtyEntityWithID(c entity.Context, id *bbq.EntityID, client *codec.PacketReadWriter) *proxyEtyEntity {

	_, err := entity.NewEntity(c, id)
	if err != nil {
		xlog.Errorln("new entity err")
		return nil
	}
	t := &proxyEtyEntity{entity: id, client: client}

	return t
}

type proxyEtyEntity struct {
	entity *bbq.EntityID

	client *codec.PacketReadWriter
}

func (t *proxyEtyEntity) RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "RegisterProxy"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterProxyResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterProxyResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyEtyEntity) SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "SyncService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(SyncServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*SyncServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxyEtyEntity) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Entity
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = t.entity
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ProxyEtyEntity
type ProxyEtyEntity interface {
	entity.IEntity

	// RegisterProxy
	RegisterProxy(c entity.Context, req *RegisterProxyRequest) (*RegisterProxyResponse, error)

	// SyncService
	SyncService(c entity.Context, req *SyncServiceRequest) (*SyncServiceResponse, error)

	// Ping
	Ping(c entity.Context, req *PingPong) (*PingPong, error)
}

func _ProxyEtyEntity_RegisterProxy_Handler(svc any, ctx entity.Context, in *RegisterProxyRequest, interceptor entity.ServerInterceptor) (*RegisterProxyResponse, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).RegisterProxy(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/RegisterProxy",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).RegisterProxy(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterProxyResponse), err

}

func _ProxyEtyEntity_RegisterProxy_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_RegisterProxy_Handler(svc, ctx, in.(*RegisterProxyRequest), interceptor)

}

func _ProxyEtyEntity_RegisterProxy_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterProxyRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_RegisterProxy_Handler(svc, ctx, in, interceptor)

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

func _ProxyEtyEntity_SyncService_Handler(svc any, ctx entity.Context, in *SyncServiceRequest, interceptor entity.ServerInterceptor) (*SyncServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).SyncService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/SyncService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).SyncService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*SyncServiceResponse), err

}

func _ProxyEtyEntity_SyncService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_SyncService_Handler(svc, ctx, in.(*SyncServiceRequest), interceptor)

}

func _ProxyEtyEntity_SyncService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(SyncServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_SyncService_Handler(svc, ctx, in, interceptor)

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

func _ProxyEtyEntity_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {
		return svc.(ProxyEtyEntity).Ping(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxyEtyEntity/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxyEtyEntity).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*PingPong), err

}

func _ProxyEtyEntity_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxyEtyEntity_Ping_Handler(svc, ctx, in.(*PingPong), interceptor)

}

func _ProxyEtyEntity_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxyEtyEntity_Ping_Handler(svc, ctx, in, interceptor)

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

var ProxyEtyEntityDesc = entity.EntityDesc{
	TypeName:    "proxypb.ProxyEtyEntity",
	HandlerType: (*ProxyEtyEntity)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterProxy": {
			MethodName:   "RegisterProxy",
			Handler:      _ProxyEtyEntity_RegisterProxy_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_RegisterProxy_Local_Handler,
		},

		"SyncService": {
			MethodName:   "SyncService",
			Handler:      _ProxyEtyEntity_SyncService_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_SyncService_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxyEtyEntity_Ping_Remote_Handler,
			LocalHandler: _ProxyEtyEntity_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}

func RegisterProxySvcService(impl ProxySvcService) {
	entity.Manager.RegisterService(&ProxySvcServiceDesc, impl)
}

func NewProxySvcServiceClient(client *codec.PacketReadWriter) *proxySvcService {
	t := &proxySvcService{client: client}
	return t
}

func NewProxySvcService(client *codec.PacketReadWriter) *proxySvcService {
	t := &proxySvcService{client: client}
	return t
}

type proxySvcService struct {
	client *codec.PacketReadWriter
}

func (t *proxySvcService) RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "RegisterInst"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterInstResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterInstResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxySvcService) RegisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "RegisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterEntityResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterEntityResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxySvcService) RegisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "RegisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxySvcService) UnregisterEntity(c entity.Context, req *RegisterEntityRequest) (*RegisterEntityResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "UnregisterEntity"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterEntityResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterEntityResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxySvcService) UnregisterService(c entity.Context, req *RegisterServiceRequest) (*RegisterServiceResponse, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "UnregisterService"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(RegisterServiceResponse)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*RegisterServiceResponse); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

func (t *proxySvcService) Ping(c entity.Context, req *PingPong) (*PingPong, error) {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = snowflake.GenUUID()
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = bbq.RequestType_RequestRequest
	pkt.Header.ServiceType = bbq.ServiceType_Service
	pkt.Header.SrcEntity = c.EntityID()
	pkt.Header.DstEntity = &bbq.EntityID{Type: "proxypb.ProxySvcService"}
	pkt.Header.Method = "Ping"
	pkt.Header.ContentType = bbq.ContentType_Proto
	pkt.Header.CompressType = bbq.CompressType_None
	pkt.Header.CheckFlags = 0
	pkt.Header.TransInfo = map[string][]byte{}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	var chanRsp chan any = make(chan any)

	err := entity.HandleCallLocalMethod(pkt, req, chanRsp)
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

		t.client.WritePacket(pkt)

		// register callback
		c.RegisterCallback(pkt.Header.RequestId, func(pkt *codec.Packet) {
			rsp := new(PingPong)
			reqbuf := pkt.PacketBody()
			err := codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp)
			if err != nil {
				chanRsp <- err
				return
			}
			chanRsp <- rsp
		})

	}

	rsp := <-chanRsp
	close(chanRsp)

	if rsp, ok := rsp.(*PingPong); ok {
		return rsp, nil
	}
	return nil, rsp.(error)

}

// ProxySvcService
type ProxySvcService interface {
	entity.IService

	// RegisterInst
	RegisterInst(c entity.Context, req *RegisterInstRequest) (*RegisterInstResponse, error)

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

func _ProxySvcService_RegisterInst_Handler(svc any, ctx entity.Context, in *RegisterInstRequest, interceptor entity.ServerInterceptor) (*RegisterInstResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).RegisterInst(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/RegisterInst",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).RegisterInst(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterInstResponse), err

}

func _ProxySvcService_RegisterInst_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_RegisterInst_Handler(svc, ctx, in.(*RegisterInstRequest), interceptor)

}

func _ProxySvcService_RegisterInst_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterInstRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_RegisterInst_Handler(svc, ctx, in, interceptor)

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

func _ProxySvcService_RegisterEntity_Handler(svc any, ctx entity.Context, in *RegisterEntityRequest, interceptor entity.ServerInterceptor) (*RegisterEntityResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).RegisterEntity(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/RegisterEntity",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).RegisterEntity(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterEntityResponse), err

}

func _ProxySvcService_RegisterEntity_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_RegisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), interceptor)

}

func _ProxySvcService_RegisterEntity_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_RegisterEntity_Handler(svc, ctx, in, interceptor)

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

func _ProxySvcService_RegisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).RegisterService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/RegisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).RegisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterServiceResponse), err

}

func _ProxySvcService_RegisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_RegisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), interceptor)

}

func _ProxySvcService_RegisterService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_RegisterService_Handler(svc, ctx, in, interceptor)

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

func _ProxySvcService_UnregisterEntity_Handler(svc any, ctx entity.Context, in *RegisterEntityRequest, interceptor entity.ServerInterceptor) (*RegisterEntityResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).UnregisterEntity(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/UnregisterEntity",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).UnregisterEntity(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterEntityResponse), err

}

func _ProxySvcService_UnregisterEntity_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_UnregisterEntity_Handler(svc, ctx, in.(*RegisterEntityRequest), interceptor)

}

func _ProxySvcService_UnregisterEntity_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterEntityRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_UnregisterEntity_Handler(svc, ctx, in, interceptor)

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

func _ProxySvcService_UnregisterService_Handler(svc any, ctx entity.Context, in *RegisterServiceRequest, interceptor entity.ServerInterceptor) (*RegisterServiceResponse, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).UnregisterService(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/UnregisterService",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).UnregisterService(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*RegisterServiceResponse), err

}

func _ProxySvcService_UnregisterService_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_UnregisterService_Handler(svc, ctx, in.(*RegisterServiceRequest), interceptor)

}

func _ProxySvcService_UnregisterService_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(RegisterServiceRequest)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_UnregisterService_Handler(svc, ctx, in, interceptor)

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

func _ProxySvcService_Ping_Handler(svc any, ctx entity.Context, in *PingPong, interceptor entity.ServerInterceptor) (*PingPong, error) {
	if interceptor == nil {
		return svc.(ProxySvcService).Ping(ctx, in)
	}

	info := &entity.ServerInfo{
		Server:     svc,
		FullMethod: "/proxypb.ProxySvcService/Ping",
	}

	handler := func(ctx entity.Context, rsp any) (any, error) {

		return svc.(ProxySvcService).Ping(ctx, in)

	}

	rsp, err := interceptor(ctx, in, info, handler)
	_ = rsp

	return rsp.(*PingPong), err

}

func _ProxySvcService_Ping_Local_Handler(svc any, ctx entity.Context, in any, interceptor entity.ServerInterceptor) (any, error) {

	return _ProxySvcService_Ping_Handler(svc, ctx, in.(*PingPong), interceptor)

}

func _ProxySvcService_Ping_Remote_Handler(svc any, ctx entity.Context, pkt *codec.Packet, interceptor entity.ServerInterceptor) {

	hdr := pkt.Header

	in := new(PingPong)
	reqbuf := pkt.PacketBody()
	err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, in)
	if err != nil {
		// nil,err
		return
	}

	rsp, err := _ProxySvcService_Ping_Handler(svc, ctx, in, interceptor)

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

var ProxySvcServiceDesc = entity.EntityDesc{
	TypeName:    "proxypb.ProxySvcService",
	HandlerType: (*ProxySvcService)(nil),
	Methods: map[string]entity.MethodDesc{

		"RegisterInst": {
			MethodName:   "RegisterInst",
			Handler:      _ProxySvcService_RegisterInst_Remote_Handler,
			LocalHandler: _ProxySvcService_RegisterInst_Local_Handler,
		},

		"RegisterEntity": {
			MethodName:   "RegisterEntity",
			Handler:      _ProxySvcService_RegisterEntity_Remote_Handler,
			LocalHandler: _ProxySvcService_RegisterEntity_Local_Handler,
		},

		"RegisterService": {
			MethodName:   "RegisterService",
			Handler:      _ProxySvcService_RegisterService_Remote_Handler,
			LocalHandler: _ProxySvcService_RegisterService_Local_Handler,
		},

		"UnregisterEntity": {
			MethodName:   "UnregisterEntity",
			Handler:      _ProxySvcService_UnregisterEntity_Remote_Handler,
			LocalHandler: _ProxySvcService_UnregisterEntity_Local_Handler,
		},

		"UnregisterService": {
			MethodName:   "UnregisterService",
			Handler:      _ProxySvcService_UnregisterService_Remote_Handler,
			LocalHandler: _ProxySvcService_UnregisterService_Local_Handler,
		},

		"Ping": {
			MethodName:   "Ping",
			Handler:      _ProxySvcService_Ping_Remote_Handler,
			LocalHandler: _ProxySvcService_Ping_Local_Handler,
		},
	},

	Metadata: "proxy.proto",
}
