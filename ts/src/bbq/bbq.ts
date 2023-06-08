import { randomUUID } from "crypto";
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

  if (!client) {
    throw "client undefined"
  }

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

      if (entityID) {
        hdr.DstEntity = entityID
      } else {
        hdr.DstEntity = EntityID.create()
        hdr.DstEntity.Type = service.typeName
      }

      hdr.ServiceType = service.serviceType;
      hdr.Method = attrs.methodName
      hdr.RequestType = RequestType.RequestRequest;

      hdr.RequestId = randomUUID()

      console.log("[sys] req:", JSON.stringify(hdr), JSON.stringify(args))

      const data = attrs.requestSerialize(args[0]);

      let resp = new Deferred<any, ERROR>()

      const rpc = client.unaryInvoke(hdr, data)

      if (attrs.responseDeserialize === undefined) {
        // is one way call
        resp.resolve(0)
      } else {
        rpc.then((ctx: UnaryContext<any>) => {
          if (!ctx.response) {
            console.log("[sys] rsp err:", ctx.error)
            resp.reject(ERROR.CLIENT_INVALID_ERR)
            return
          }
          if (ctx.response.Header.ErrCode) {
            console.log("[sys] rsp err:", ctx.response.Header.ErrCode)
            resp.reject(ctx.response.Header.ErrCode)
            return
          }

          let rsp = attrs.responseDeserialize(ctx.response?.PacketBody())
          // console.log("ctx resp:", sayResp)
          console.log("[sys] rsp:", JSON.stringify(ctx.response.Header), JSON.stringify(rsp))
          resp.resolve({ error: undefined, response: rsp })
        });
      }

      return resp.promise
    }

    ServiceClientImpl.prototype[name] = methodFunc;
  })

  return new ServiceClientImpl()
}


