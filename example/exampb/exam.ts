/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "exampb";

export interface SayHelloRequest {
  text: string;
  CLientID: Long;
}

export interface SayHelloResponse {
  text: string;
}

function createBaseSayHelloRequest(): SayHelloRequest {
  return { text: "", CLientID: Long.UZERO };
}

export const SayHelloRequest = {
  encode(message: SayHelloRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    if (!message.CLientID.isZero()) {
      writer.uint32(16).uint64(message.CLientID);
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
          message.CLientID = reader.uint64() as Long;
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
      CLientID: isSet(object.CLientID) ? Long.fromValue(object.CLientID) : Long.UZERO,
    };
  },

  toJSON(message: SayHelloRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.CLientID !== undefined && (obj.CLientID = (message.CLientID || Long.UZERO).toString());
    return obj;
  },

  create(base?: DeepPartial<SayHelloRequest>): SayHelloRequest {
    return SayHelloRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SayHelloRequest>): SayHelloRequest {
    const message = createBaseSayHelloRequest();
    message.text = object.text ?? "";
    message.CLientID = (object.CLientID !== undefined && object.CLientID !== null)
      ? Long.fromValue(object.CLientID)
      : Long.UZERO;
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

/** 客户端 */
export type NoRespDefinition = typeof NoRespDefinition;
export const NoRespDefinition = {
  name: "NoResp",
  fullName: "exampb.NoResp",
  methods: {
    sayHello: {
      name: "SayHello",
      requestType: SayHelloRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Long ? string | number | Long : T extends Array<infer U> ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
