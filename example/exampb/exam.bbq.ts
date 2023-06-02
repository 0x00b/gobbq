import { SayHelloRequest, SayHelloResponse } from "./exam"

import { ServiceType } from "./bbq";
import { UnaryResponse } from "../../ts/src/context/unary";
import { Client } from "../../ts/src";
import { makeClientConstructor } from "../../ts/src/bbq/bbq";
import { EntityID } from "../../proto/bbq/bbq";

export type EchoDefinition = typeof EchoDefinition;
export const EchoDefinition = {
  typeName: "exampb.Echo",
  serviceType: ServiceType.Service,//or Entity
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: SayHelloResponse,
      requestSerialize: serialize_exampb_SayHelloRequest,
      requestDeserialize: deserialize_exampb_SayHelloRequest,
      responseSerialize: serialize_exampb_SayHelloResponse,
      responseDeserialize: deserialize_exampb_SayHelloResponse,
    },
  },
} as const;

function serialize_exampb_SayHelloRequest(req: SayHelloRequest): Buffer {
  return Buffer.from(SayHelloRequest.encode(req).finish())
}

function deserialize_exampb_SayHelloRequest(input: Uint8Array): SayHelloRequest {
  return SayHelloRequest.decode(input)
}
function serialize_exampb_SayHelloResponse(req: SayHelloResponse): Buffer {
  return Buffer.from(SayHelloResponse.encode(req).finish())
}

function deserialize_exampb_SayHelloResponse(input: Uint8Array): SayHelloResponse {
  return SayHelloRequest.decode(input)
}

export interface EchoService {
  SayHello(request: SayHelloRequest): UnaryResponse<SayHelloResponse>;
}

export function NewEchoService(client: Client<any>): EchoService {
  return makeClientConstructor(client, EchoDefinition) as unknown as EchoService
}

export function NewEchoEntity(client: Client<any>, entityID: EntityID): EchoService {
  return makeClientConstructor(client, EchoDefinition, entityID) as unknown as EchoService
}