import {SayHelloRequest, SayHelloResponse} from "./exam"
import {CompressType,ServiceType,RequestType, ContentType}from  "./bbq"

class EchoService {
    async SayHello(c: EntityContext, req: SayHelloRequest): Promise<SayHelloResponse> {
        const pkt = codec.NewPacket();
        const release = pkt.release.bind(pkt);

        try {
            pkt.Header.Version = 1;
            pkt.Header.RequestId = snowflake.GenUUID();
            pkt.Header.Timeout = 10;
            pkt.Header.RequestType = RequestType.RequestRequest;
            pkt.Header.ServiceType = ServiceType.Service;
            pkt.Header.SrcEntity = c.EntityID();
            pkt.Header.DstEntity = { Type: "exampb.EchoService" };
            pkt.Header.Method = "SayHello";
            pkt.Header.ContentType =  ContentType.Proto;
            pkt.Header.CompressType = CompressType.None;
            pkt.Header.CheckFlags = 0;
            pkt.Header.TransInfo = {};
            pkt.Header.ErrCode = 0;
            pkt.Header.ErrMsg = "";

            const chanRsp: any = new Promise((resolve, reject) => {
                const callback = (pkt: codec.Packet) => {
                    const rsp = new SayHelloResponse();
                    const reqbuf = pkt.PacketBody();
                    const err = codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp);
                    if (err != null) {
                        reject(err);
                        return;
                    }
                    resolve(rsp);
                };
                entity.RegisterCallback(c, pkt.Header.RequestId, callback);
            });

            const etyMgr = entity.GetEntityMgr(c);
            if (etyMgr == null) {
                throw new Error("bad context");
            }
            const err = await etyMgr.LocalCall(pkt, req, chanRsp);
            if (err != null) {
                if (!entity.NotMyMethod(err)) {
                    throw err;
                }

                const hdrBytes = await codec.GetCodec(bbq.ContentType_Proto).Marshal(req);
                pkt.WriteBody(hdrBytes);

                // register callback first, than SendPackt
                entity.RegisterCallback(c, pkt.Header.RequestId, (pkt: codec.Packet) => {
                    const rsp = new SayHelloResponse();
                    const reqbuf = pkt.PacketBody();
                    const err = codec.GetCodec(pkt.Header.GetContentType()).Unmarshal(reqbuf, rsp);
                    if (err != null) {
                        chanRsp.reject(err);
                        return;
                    }
                    chanRsp.resolve(rsp);
                });

                const remoteEntityManager = entity.GetRemoteEntityManager(c);
                await remoteEntityManager.SendPackt(pkt);
            }

            const rsp = await chanRsp;

            if (rsp instanceof SayHelloResponse) {
                return rsp;
            }
            throw rsp;
        } catch (e) {
            throw e;
        } finally {
            release();
        }
    }
}
// EchoService
interface EchoService {
    iservice: entity.IService;
  
    // SayHello
    sayHello(ctx: entity.Context, request: SayHelloRequest): SayHelloResponse;
  }
  
  function _EchoService_SayHello_Handler(svc: any, ctx: entity.Context, req: SayHelloRequest, interceptor: Entity.ServerInterceptor): SayHelloResponse {
    if (!interceptor) {
      return svc.echoService.sayHello(ctx, req);
    }
  
    const info = {
      server: svc,
      fullMethod: "/exampb.EchoService/SayHello";
    };
  
    const handler = function(ctx: entity.Context, res: any) {
      return svc.echoService.sayHello(ctx, req);
    };
  
    const res = interceptor(ctx, req, info, handler);
  
    return res as SayHelloResponse;
  }
  
  function _EchoService_SayHello_Local_Handler(svc: any, ctx: Entity.Context, req: any, interceptor: Entity.ServerInterceptor): any {
    return _EchoService_SayHello_Handler(svc, ctx, req as SayHelloRequest, interceptor);
  }
  
  function _EchoService_SayHello_Remote_Handler(svc: any, ctx: Entity.Context, pkt: codec.Packet, interceptor: Entity.ServerInterceptor) {
    const hdr = pkt.header;
  
    const req = new SayHelloRequest();
    const reqBuf = pkt.packetBody();
    const err = codec.getCodec(hdr.getContentType()).unmarshal(reqBuf, req);
  
    if (err) {
      // nil, err
      return;
    }
  
    const res = _EchoService_SayHello_Handler(svc, ctx, req, interceptor);
  
    const npkt = codec.newPacket();
    
    npkt.header.version = hdr.version;
    npkt.header.requestId = hdr.requestId;
    npkt.header.timeout = hdr.timeout;
    npkt.header.requestType