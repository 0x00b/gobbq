import { SayHelloRequest, SayHelloResponse } from "../../../example/exampb/exam";
import { Echo, EchoDefinition } from "../../../example/exampb/exam.bbq";
import { EntityID, Header, RequestType, ServiceType } from "../../../proto/bbq/bbq";
import { Client } from "../client";
import { UnaryContext } from "../context";
import { UnaryResponse } from "../context/unary";
import { Context } from "../dispatcher/context";
import { ServiceDefinition } from "../dispatcher/service";
import { ERROR } from "../error";
import { Deferred } from "../utils";

export function makeClientConstructor(
  client: Client<any>,
  service: ServiceDefinition,
  entityID?: EntityID
) {

  if (service.serviceType == ServiceType.Entity && !entityID) {
    throw "need entity id"
  }

  class ServiceClientImpl {
    [methodName: string]: Function;
  }

  let methods = service.methods
  Object.keys(methods).forEach((name) => {

    let attrs = methods[name]

    const methodFunc = (...args: any[]) => {

      let hdr = Header.create()

      hdr.DstEntity = EntityID.create()
      hdr.DstEntity.Type = service.typeName

      hdr.ServiceType = service.serviceType;
      hdr.Method = attrs.methodName

      hdr.RequestType = RequestType.RequestRequest;

      hdr.RequestId = "1111"
      hdr.Timeout = 1000

      //    hdr.SrcEntity = c.ID;

      console.log(JSON.stringify(args))

      const data = attrs.requestSerialize(args[0]);

      let resp = new Deferred<any, ERROR>()

      const rpc = client.unaryInvoke(hdr, data).then((ctx: UnaryContext<any>) => {
        if (!ctx.response) {
          console.log("ctx resp is undefined")
          resp.reject(ERROR.CLIENT_CANCELED_ERR)
          return
        }
        if (ctx.response.Header.ErrCode) {
          resp.reject(ctx.response.Header.ErrCode)
          return
        }

        let sayResp = attrs.responseDeserialize(ctx.response?.PacketBody())
        console.log("ctx resp:", sayResp)
        resp.resolve(sayResp)
      });

      return resp.promise
    }

    ServiceClientImpl.prototype[name] = methodFunc;
  })

  let x = new ServiceClientImpl()

  console.log(JSON.stringify(x))

  return x;
}


class EchoImpl {
  SayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {
    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse:", request.text)

    let rsp = SayHelloResponse.create()
    rsp.text = "xxxx"
    return rsp
  }
}

function test() {

  const remote = {
    port: 8899,
    host: 'localhost',
    protocol: 'kcp',
  } as const;

  const impl: any = new EchoImpl();
  let client = new Client(EchoDefinition, impl, { remote })

  let c = makeClientConstructor(client, EchoDefinition) as unknown as Echo

  let rsp = c.SayHello({ text: "request", CLientID: undefined })

  console.log("say resp 11", rsp)

  rsp.then((rsp) => {
    if (rsp instanceof Error) {
      console.log("error", rsp)
      return
    }

    console.log("succ rsp:", rsp)
  })

  let rsp2 = c.SayHello({ text: "request", CLientID: undefined })

  console.log("say resp22", rsp2)

}

test()
