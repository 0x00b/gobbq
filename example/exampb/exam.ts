/* eslint-disable */
import _m0 from "../../third/node_modules/protobufjs/minimal";
import { EntityID } from "./bbq";

export const protobufPackage = "exampb";

export interface SayHelloRequest {
  text: string;
  CLientID: EntityID | undefined;
}

export interface SayHelloResponse {
  text: string;
}

function createBaseSayHelloRequest(): SayHelloRequest {
  return { text: "", CLientID: undefined };
}

export const SayHelloRequest = {
  encode(message: SayHelloRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    if (message.CLientID !== undefined) {
      EntityID.encode(message.CLientID, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SayHelloRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSayHelloRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        case 2:
          message.CLientID = EntityID.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SayHelloRequest {
    return {
      text: isSet(object.text) ? String(object.text) : "",
      CLientID: isSet(object.CLientID) ? EntityID.fromJSON(object.CLientID) : undefined,
    };
  },

  toJSON(message: SayHelloRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.CLientID !== undefined && (obj.CLientID = message.CLientID ? EntityID.toJSON(message.CLientID) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SayHelloRequest>): SayHelloRequest {
    return SayHelloRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SayHelloRequest>): SayHelloRequest {
    const message = createBaseSayHelloRequest();
    message.text = object.text ?? "";
    message.CLientID = (object.CLientID !== undefined && object.CLientID !== null)
      ? EntityID.fromPartial(object.CLientID)
      : undefined;
    return message;
  },
};

function createBaseSayHelloResponse(): SayHelloResponse {
  return { text: "" };
}

export const SayHelloResponse = {
  encode(message: SayHelloResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SayHelloResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSayHelloResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SayHelloResponse {
    return { text: isSet(object.text) ? String(object.text) : "" };
  },

  toJSON(message: SayHelloResponse): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    return obj;
  },

  create(base?: DeepPartial<SayHelloResponse>): SayHelloResponse {
    return SayHelloResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SayHelloResponse>): SayHelloResponse {
    const message = createBaseSayHelloResponse();
    message.text = object.text ?? "";
    return message;
  },
};

export type EchoDefinition = typeof EchoDefinition;
export const EchoDefinition = {
  name: "Echo",
  fullName: "exampb.Echo",
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

export type EchoEtyDefinition = typeof EchoEtyDefinition;
export const EchoEtyDefinition = {
  name: "EchoEty",
  fullName: "exampb.EchoEty",
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

export type EchoSvc2Definition = typeof EchoSvc2Definition;
export const EchoSvc2Definition = {
  name: "EchoSvc2",
  fullName: "exampb.EchoSvc2",
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
