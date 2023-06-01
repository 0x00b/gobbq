import { EntityID, Header, RequestType, ServiceType } from "../../../proto/bbq/bbq";
import { Client } from "../client";
import { UnaryContext } from "../context";
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


