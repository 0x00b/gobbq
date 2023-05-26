

 import { client } from '../src';
 import * as bbq from '../../proto/bbq/bbq';
  
 
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
 
   let hdr =  bbq.Header.create()
   hdr.RequestId = "1111"
   hdr.Timeout=1000
   hdr.RequestType = bbq.RequestType.RequestRequest;

   hdr.ServiceType = bbq.ServiceType.Service;
   hdr.Method = "SayHello";

//    hdr.SrcEntity = c.ID;
//    hdr.DstEntity = dstEntity

   const data = Buffer.from(bbq.Header.encode(hdr).finish());

   const contentType = bbq.ContentType.Proto;
 
   /**
    * 远程接入点；
    * 作为示例，我们直接传入 mock server 的接入点；
    * 实际使用中，一般由名字服务中间件提供接入点；
    */
   const remote = {
     port: 8899,
     host: 'localhost',
     protocol: 'tcp',
   } as const;
   let timeout = hdr.Timeout
   const rpc = await client.channel.unaryInvoke(hdr, data, { contentType, remote, timeout });
 
   // 调用错误
   console.log(__filename, rpc.error);
   //  请求消息
   console.log(__filename, rpc.request);
   //  响应消息
   console.log(__filename, rpc.response);
 }
 
 invoke();
 