import { SayHelloRequest, SayHelloResponse } from "./exam"
import { v4 as uuidv4 } from 'uuid';
import * as packet from "../../ts/codec/packet"
import * as client from "../../ts/client/client"
import * as bbq from "../../proto/bbq/bbq"

export class EchoService {
    async SayHello(c: client.Client, req: SayHelloRequest): Promise<SayHelloResponse> {
        const pkt = new packet.Packet();

        try {
            let dstEntity = bbq.EntityID.create()
            dstEntity.Type = "exampb.EchoService"

            pkt.Header.Version = 1;
            pkt.Header.RequestId = uuidv4();
            pkt.Header.Timeout = 10;
            pkt.Header.RequestType = bbq.RequestType.RequestRequest;
            pkt.Header.ServiceType = bbq.ServiceType.Service;
            pkt.Header.SrcEntity = c.ID;
            pkt.Header.DstEntity = dstEntity
            pkt.Header.Method = "SayHello";
            pkt.Header.ContentType = bbq.ContentType.Proto;
            pkt.Header.CompressType = bbq.CompressType.None;
            pkt.Header.CheckFlags = 0;
            pkt.Header.TransInfo = {};
            pkt.Header.ErrCode = 0;
            pkt.Header.ErrMsg = "";

            const chanRsp: any = new Promise((resolve, reject) => {
                const callback = (pkt: packet.Packet) => {
                    const reqbuf = pkt.PacketBody();
                    if (reqbuf != null) {
                        const rsp = SayHelloResponse.decode(reqbuf);
                        resolve(rsp);
                    }
                    reject()
                };
                c.RegisterCallback(pkt.Header.RequestId, callback);
            });

            const bodyBytes = SayHelloRequest.encode(req).finish();
            pkt.WriteBody(Buffer.from(bodyBytes));

            c.SendPacket(pkt);

            const rsp = await chanRsp;

            // if (rsp instanceof SayHelloResponse) {
            //     return rsp;
            // }
            throw rsp;
        } catch (e) {
            throw e;
        } finally {
        }
    }
}

// EchoService
interface IEchoService {

    // SayHello
    SayHello(request: SayHelloRequest): SayHelloResponse;
}

function _EchoService_SayHello_Remote_Handler(svc: IEchoService, pkt: packet.Packet,) {

    const reqBuf = pkt.PacketBody();
    if (reqBuf == null) {
        return
    }

    const req = SayHelloRequest.decode(Buffer.from(reqBuf));

    // const res = _EchoService_SayHello_Handler(svc, ctx, req, interce
    let rsp = svc.SayHello(req);

    const npkt = new packet.Packet();

    let hdr = pkt.Header

    let sstEntity = bbq.EntityID.create()
    sstEntity.Type = "exampb.EchoService"

    npkt.Header.Version = hdr.Version;
    npkt.Header.RequestId = hdr.RequestId;
    npkt.Header.Timeout = hdr.Timeout;
    npkt.Header.RequestType = bbq.RequestType.RequestRespone;
    npkt.Header.ServiceType = bbq.ServiceType.Service;
    npkt.Header.SrcEntity = sstEntity;
    npkt.Header.DstEntity = hdr.SrcEntity
    npkt.Header.Method = "SayHello";
    npkt.Header.ContentType = bbq.ContentType.Proto;
    npkt.Header.CompressType = bbq.CompressType.None;
    npkt.Header.CheckFlags = 0;
    npkt.Header.TransInfo = {};
    npkt.Header.ErrCode = 0;
    npkt.Header.ErrMsg = "";

    const rspBuff = SayHelloResponse.encode(rsp).finish();

    npkt.WriteBody(Buffer.from(rspBuff))
}


/** 客户端 */
export type ClientDefinition = typeof ClientDefinition;
export const ClientDefinition = {
    name: "Client",
    fullName: "exampb.Client",
    methods: {
        sayHello: {
            name: "SayHello",
            requestType: SayHelloRequest,
            requestStream: false,
            responseType: SayHelloResponse,
            responseStream: false,
            options: {},
        },
    },
} as const;
