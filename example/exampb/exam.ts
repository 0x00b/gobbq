/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
