/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { EntityID } from "./bbq";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "gatepb";

export interface PingPong {
}

export interface RegisterClientRequest {
  EntityID: EntityID | undefined;
}

export interface RegisterClientResponse {
  EntityID: EntityID | undefined;
}

function createBasePingPong(): PingPong {
  return {};
}

export const PingPong = {
  encode(_: PingPong, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PingPong {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePingPong();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): PingPong {
    return {};
  },

  toJSON(_: PingPong): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<PingPong>): PingPong {
    return PingPong.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<PingPong>): PingPong {
    const message = createBasePingPong();
    return message;
  },
};

function createBaseRegisterClientRequest(): RegisterClientRequest {
  return { EntityID: undefined };
}

export const RegisterClientRequest = {
  encode(message: RegisterClientRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.EntityID !== undefined) {
      EntityID.encode(message.EntityID, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterClientRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterClientRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.EntityID = EntityID.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RegisterClientRequest {
    return { EntityID: isSet(object.EntityID) ? EntityID.fromJSON(object.EntityID) : undefined };
  },

  toJSON(message: RegisterClientRequest): unknown {
    const obj: any = {};
    message.EntityID !== undefined && (obj.EntityID = message.EntityID ? EntityID.toJSON(message.EntityID) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RegisterClientRequest>): RegisterClientRequest {
    return RegisterClientRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RegisterClientRequest>): RegisterClientRequest {
    const message = createBaseRegisterClientRequest();
    message.EntityID = (object.EntityID !== undefined && object.EntityID !== null)
      ? EntityID.fromPartial(object.EntityID)
      : undefined;
    return message;
  },
};

function createBaseRegisterClientResponse(): RegisterClientResponse {
  return { EntityID: undefined };
}

export const RegisterClientResponse = {
  encode(message: RegisterClientResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.EntityID !== undefined) {
      EntityID.encode(message.EntityID, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterClientResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterClientResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.EntityID = EntityID.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RegisterClientResponse {
    return { EntityID: isSet(object.EntityID) ? EntityID.fromJSON(object.EntityID) : undefined };
  },

  toJSON(message: RegisterClientResponse): unknown {
    const obj: any = {};
    message.EntityID !== undefined && (obj.EntityID = message.EntityID ? EntityID.toJSON(message.EntityID) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RegisterClientResponse>): RegisterClientResponse {
    return RegisterClientResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RegisterClientResponse>): RegisterClientResponse {
    const message = createBaseRegisterClientResponse();
    message.EntityID = (object.EntityID !== undefined && object.EntityID !== null)
      ? EntityID.fromPartial(object.EntityID)
      : undefined;
    return message;
  },
};

export interface Gate {
  RegisterClient(request: RegisterClientRequest): Promise<RegisterClientResponse>;
  UnregisterClient(request: RegisterClientRequest): Promise<Empty>;
  Ping(request: PingPong): Promise<PingPong>;
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
