/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "exampb";

export enum Enum {
  XXX = 0,
  XXX2 = 1,
  UNRECOGNIZED = -1,
}

export function enumFromJSON(object: any): Enum {
  switch (object) {
    case 0:
    case "XXX":
      return Enum.XXX;
    case 1:
    case "XXX2":
      return Enum.XXX2;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Enum.UNRECOGNIZED;
  }
}

export function enumToJSON(object: Enum): string {
  switch (object) {
    case Enum.XXX:
      return "XXX";
    case Enum.XXX2:
      return "XXX2";
    case Enum.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SayHelloRequest {
  text: string;
  CLientID: Long;
}

export interface SayHelloResponse {
  text: string;
}

export interface EchoProperty {
  Text: string;
  test:
    | SayHelloRequest
    | undefined;
  /** 先不支持import类型 */
  test2: Empty[];
  test3: { [key: number]: string };
  test4: Uint8Array;
  test5: Long;
  test6: number;
  TEST7: Enum;
  test8: number;
  test9: number[];
}

export interface EchoProperty_Test3Entry {
  key: number;
  value: string;
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

function createBaseEchoProperty(): EchoProperty {
  return {
    Text: "",
    test: undefined,
    test2: [],
    test3: {},
    test4: new Uint8Array(),
    test5: Long.ZERO,
    test6: 0,
    TEST7: 0,
    test8: 0,
    test9: [],
  };
}

export const EchoProperty = {
  encode(message: EchoProperty, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Text !== "") {
      writer.uint32(10).string(message.Text);
    }
    if (message.test !== undefined) {
      SayHelloRequest.encode(message.test, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.test2) {
      Empty.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    Object.entries(message.test3).forEach(([key, value]) => {
      EchoProperty_Test3Entry.encode({ key: key as any, value }, writer.uint32(34).fork()).ldelim();
    });
    if (message.test4.length !== 0) {
      writer.uint32(42).bytes(message.test4);
    }
    if (!message.test5.isZero()) {
      writer.uint32(49).sfixed64(message.test5);
    }
    if (message.test6 !== 0) {
      writer.uint32(61).sfixed32(message.test6);
    }
    if (message.TEST7 !== 0) {
      writer.uint32(64).int32(message.TEST7);
    }
    if (message.test8 !== 0) {
      writer.uint32(72).sint32(message.test8);
    }
    writer.uint32(82).fork();
    for (const v of message.test9) {
      writer.sint32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EchoProperty {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEchoProperty();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Text = reader.string();
          break;
        case 2:
          message.test = SayHelloRequest.decode(reader, reader.uint32());
          break;
        case 3:
          message.test2.push(Empty.decode(reader, reader.uint32()));
          break;
        case 4:
          const entry4 = EchoProperty_Test3Entry.decode(reader, reader.uint32());
          if (entry4.value !== undefined) {
            message.test3[entry4.key] = entry4.value;
          }
          break;
        case 5:
          message.test4 = reader.bytes();
          break;
        case 6:
          message.test5 = reader.sfixed64() as Long;
          break;
        case 7:
          message.test6 = reader.sfixed32();
          break;
        case 8:
          message.TEST7 = reader.int32() as any;
          break;
        case 9:
          message.test8 = reader.sint32();
          break;
        case 10:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.test9.push(reader.sint32());
            }
          } else {
            message.test9.push(reader.sint32());
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EchoProperty {
    return {
      Text: isSet(object.Text) ? String(object.Text) : "",
      test: isSet(object.test) ? SayHelloRequest.fromJSON(object.test) : undefined,
      test2: Array.isArray(object?.test2) ? object.test2.map((e: any) => Empty.fromJSON(e)) : [],
      test3: isObject(object.test3)
        ? Object.entries(object.test3).reduce<{ [key: number]: string }>((acc, [key, value]) => {
          acc[Number(key)] = String(value);
          return acc;
        }, {})
        : {},
      test4: isSet(object.test4) ? bytesFromBase64(object.test4) : new Uint8Array(),
      test5: isSet(object.test5) ? Long.fromValue(object.test5) : Long.ZERO,
      test6: isSet(object.test6) ? Number(object.test6) : 0,
      TEST7: isSet(object.TEST7) ? enumFromJSON(object.TEST7) : 0,
      test8: isSet(object.test8) ? Number(object.test8) : 0,
      test9: Array.isArray(object?.test9) ? object.test9.map((e: any) => Number(e)) : [],
    };
  },

  toJSON(message: EchoProperty): unknown {
    const obj: any = {};
    message.Text !== undefined && (obj.Text = message.Text);
    message.test !== undefined && (obj.test = message.test ? SayHelloRequest.toJSON(message.test) : undefined);
    if (message.test2) {
      obj.test2 = message.test2.map((e) => e ? Empty.toJSON(e) : undefined);
    } else {
      obj.test2 = [];
    }
    obj.test3 = {};
    if (message.test3) {
      Object.entries(message.test3).forEach(([k, v]) => {
        obj.test3[k] = v;
      });
    }
    message.test4 !== undefined &&
      (obj.test4 = base64FromBytes(message.test4 !== undefined ? message.test4 : new Uint8Array()));
    message.test5 !== undefined && (obj.test5 = (message.test5 || Long.ZERO).toString());
    message.test6 !== undefined && (obj.test6 = Math.round(message.test6));
    message.TEST7 !== undefined && (obj.TEST7 = enumToJSON(message.TEST7));
    message.test8 !== undefined && (obj.test8 = Math.round(message.test8));
    if (message.test9) {
      obj.test9 = message.test9.map((e) => Math.round(e));
    } else {
      obj.test9 = [];
    }
    return obj;
  },

  create(base?: DeepPartial<EchoProperty>): EchoProperty {
    return EchoProperty.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EchoProperty>): EchoProperty {
    const message = createBaseEchoProperty();
    message.Text = object.Text ?? "";
    message.test = (object.test !== undefined && object.test !== null)
      ? SayHelloRequest.fromPartial(object.test)
      : undefined;
    message.test2 = object.test2?.map((e) => Empty.fromPartial(e)) || [];
    message.test3 = Object.entries(object.test3 ?? {}).reduce<{ [key: number]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[Number(key)] = String(value);
      }
      return acc;
    }, {});
    message.test4 = object.test4 ?? new Uint8Array();
    message.test5 = (object.test5 !== undefined && object.test5 !== null) ? Long.fromValue(object.test5) : Long.ZERO;
    message.test6 = object.test6 ?? 0;
    message.TEST7 = object.TEST7 ?? 0;
    message.test8 = object.test8 ?? 0;
    message.test9 = object.test9?.map((e) => e) || [];
    return message;
  },
};

function createBaseEchoProperty_Test3Entry(): EchoProperty_Test3Entry {
  return { key: 0, value: "" };
}

export const EchoProperty_Test3Entry = {
  encode(message: EchoProperty_Test3Entry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).int32(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EchoProperty_Test3Entry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEchoProperty_Test3Entry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.int32();
          break;
        case 2:
          message.value = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EchoProperty_Test3Entry {
    return { key: isSet(object.key) ? Number(object.key) : 0, value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: EchoProperty_Test3Entry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = Math.round(message.key));
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<EchoProperty_Test3Entry>): EchoProperty_Test3Entry {
    return EchoProperty_Test3Entry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EchoProperty_Test3Entry>): EchoProperty_Test3Entry {
    const message = createBaseEchoProperty_Test3Entry();
    message.key = object.key ?? 0;
    message.value = object.value ?? "";
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

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
