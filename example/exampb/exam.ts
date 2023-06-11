/* eslint-disable */
import * as Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "exampb";

export interface SayHelloRequest {
  text: string;
  CLientID: number;
}

export interface SayHelloResponse {
  text: string;
}

function createBaseSayHelloRequest(): SayHelloRequest {
  return { text: "", CLientID: 0 };
}

export const SayHelloRequest = {
  encode(message: SayHelloRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    if (message.CLientID !== 0) {
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
          message.CLientID = longToNumber(reader.uint64() as Long);
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
      CLientID: isSet(object.CLientID) ? Number(object.CLientID) : 0,
    };
  },

  toJSON(message: SayHelloRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.CLientID !== undefined && (obj.CLientID = Math.round(message.CLientID));
    return obj;
  },

  create(base?: DeepPartial<SayHelloRequest>): SayHelloRequest {
    return SayHelloRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SayHelloRequest>): SayHelloRequest {
    const message = createBaseSayHelloRequest();
    message.text = object.text ?? "";
    message.CLientID = object.CLientID ?? 0;
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

export interface Echo {
  SayHello(request: SayHelloRequest): Promise<SayHelloResponse>;
}

export interface EchoEty {
  SayHello(request: SayHelloRequest): Promise<SayHelloResponse>;
}

export interface EchoSvc2 {
  SayHello(request: SayHelloRequest): Promise<SayHelloResponse>;
}

/** 客户端 */
export interface Client {
  SayHello(request: SayHelloRequest): Promise<SayHelloResponse>;
}

/** 客户端 */
export interface NoResp {
  SayHello(request: SayHelloRequest): Promise<Empty>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
