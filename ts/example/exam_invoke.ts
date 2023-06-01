

import { SayHelloRequest, SayHelloResponse } from '../../example/exampb/exam';
import { EchoDefinition } from '../../example/exampb/exam.bbq';
import * as bbq from '../../proto/bbq/bbq';
import { Client } from '../src';
import { Context } from '../src/dispatcher/context';


class Echo {
  SayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {
    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse:", request.text)

    let rsp = SayHelloResponse.create()
    rsp.text = "xxxx"
    return rsp
  }
}



// var ClientService = exports.ClientService = {
//     sayHello: {
//       Method: '/exampb.Client/SayHello',
//       ServiceType: bbq.RequestType.RequestRequest,
//       requestStream: false,
//       responseStream: false,
//       requestType: exam_pb.SayHelloRequest,
//       responseType: exam_pb.SayHelloResponse,
//     },
//   };

async function invoke() {
  /* 接口名 func */

  let hdr = bbq.Header.create()
  hdr.DstEntity = bbq.EntityID.create()

  hdr.ServiceType = EchoDefinition.serviceType;
  hdr.DstEntity.Type = EchoDefinition.typeName
  hdr.Method = EchoDefinition.methods.SayHello.methodName

  hdr.RequestType = bbq.RequestType.RequestRequest;

  hdr.RequestId = "1111"
  hdr.Timeout = 1000

  //    hdr.SrcEntity = c.ID;

  let req = SayHelloRequest.create()
  req.text = "request"
  const data = Buffer.from(SayHelloRequest.encode(req).finish());

  const contentType = bbq.ContentType.Proto;

  /**
   * 远程接入点；
   * 作为示例，我们直接传入 mock server 的接入点；
   * 实际使用中，一般由名字服务中间件提供接入点；
   */
  const remote = {
    port: 8899,
    host: 'localhost',
    protocol: 'kcp',
  } as const;
  let timeout = hdr.Timeout

  const impl: any = new Echo();
  let client = new Client(EchoDefinition, impl, { remote })

  const rpc = await client.unaryInvoke(hdr, data, { contentType, remote, timeout });

  // 调用错误
  console.log("err", rpc.error);
  //  请求消息
  console.log("requst", rpc.request);
  //  响应消息
  console.log("response", rpc.response);
}

invoke();

