/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bbq";

export enum ContentType {
  Proto = 0,
  UNRECOGNIZED = -1,
}

export function contentTypeFromJSON(object: any): ContentType {
  switch (object) {
    case 0:
    case "Proto":
      return ContentType.Proto;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ContentType.UNRECOGNIZED;
  }
}

export function contentTypeToJSON(object: ContentType): string {
  switch (object) {
    case ContentType.Proto:
      return "Proto";
    case ContentType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum CompressType {
  None = 0,
  Gzip = 1,
  UNRECOGNIZED = -1,
}

export function compressTypeFromJSON(object: any): CompressType {
  switch (object) {
    case 0:
    case "None":
      return CompressType.None;
    case 1:
    case "Gzip":
      return CompressType.Gzip;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CompressType.UNRECOGNIZED;
  }
}

export function compressTypeToJSON(object: CompressType): string {
  switch (object) {
    case CompressType.None:
      return "None";
    case CompressType.Gzip:
      return "Gzip";
    case CompressType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum RequestType {
  RequestRequest = 0,
  RequestRespone = 1,
  UNRECOGNIZED = -1,
}

export function requestTypeFromJSON(object: any): RequestType {
  switch (object) {
    case 0:
    case "RequestRequest":
      return RequestType.RequestRequest;
    case 1:
    case "RequestRespone":
      return RequestType.RequestRespone;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RequestType.UNRECOGNIZED;
  }
}

export function requestTypeToJSON(object: RequestType): string {
  switch (object) {
    case RequestType.RequestRequest:
      return "RequestRequest";
    case RequestType.RequestRespone:
      return "RequestRespone";
    case RequestType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ServiceType {
  /** Entity - 请求entity，需要提供entity id， entity是有ID的service, entity可以创建很多 */
  Entity = 0,
  /** Service - 请求service，只需要提供完整接口名，service是单例entity，只能有一个 */
  Service = 1,
  UNRECOGNIZED = -1,
}

export function serviceTypeFromJSON(object: any): ServiceType {
  switch (object) {
    case 0:
    case "Entity":
      return ServiceType.Entity;
    case 1:
    case "Service":
      return ServiceType.Service;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ServiceType.UNRECOGNIZED;
  }
}

export function serviceTypeToJSON(object: ServiceType): string {
  switch (object) {
    case ServiceType.Entity:
      return "Entity";
    case ServiceType.Service:
      return "Service";
    case ServiceType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface EntityID {
  /** 记录Entity所在的proxy */
  ProxyID: string;
  /** proxy上记录有哪些InstID */
  InstID: string;
  /** 具体的Entity在Inst上 */
  ID: string;
  /** 具体的entity的类型 */
  Type: string;
}

/** 请求协议头 */
export interface Header {
  /** 协议版本 */
  Version: number;
  /** 请求唯一id */
  RequestId: string;
  /** 请求的超时时间，单位ms */
  Timeout: number;
  /** 是请求包，还是返回包 */
  RequestType: RequestType;
  /** sverice or entity */
  ServiceType: ServiceType;
  /** 调用的原EntityID */
  SrcEntity:
    | EntityID
    | undefined;
  /** 调用的目的EntityID */
  DstEntity:
    | EntityID
    | undefined;
  /** 规范格式: 接口名 */
  Method: string;
  /**
   * 请求数据的序列化类型
   * 比如: proto/jce/json, 默认proto
   * 具体值与ContentEncodeType对应
   */
  ContentType: ContentType;
  /**
   * 请求数据使用的压缩方式
   * 比如: gzip/snappy/..., 默认不使用
   * 具体值与CompressType对应
   */
  CompressType: CompressType;
  /** 是否检查包是否正确 */
  CheckFlags: number;
  /** 附加信息 */
  TransInfo: { [key: string]: Uint8Array };
  /** 返回值 */
  ErrCode: number;
  ErrMsg: string;
}

export interface Header_TransInfoEntry {
  key: string;
  value: Uint8Array;
}

function createBaseEntityID(): EntityID {
  return { ProxyID: "", InstID: "", ID: "", Type: "" };
}

export const EntityID = {
  encode(message: EntityID, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ProxyID !== "") {
      writer.uint32(10).string(message.ProxyID);
    }
    if (message.InstID !== "") {
      writer.uint32(18).string(message.InstID);
    }
    if (message.ID !== "") {
      writer.uint32(26).string(message.ID);
    }
    if (message.Type !== "") {
      writer.uint32(34).string(message.Type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EntityID {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntityID();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ProxyID = reader.string();
          break;
        case 2:
          message.InstID = reader.string();
          break;
        case 3:
          message.ID = reader.string();
          break;
        case 4:
          message.Type = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EntityID {
    return {
      ProxyID: isSet(object.ProxyID) ? String(object.ProxyID) : "",
      InstID: isSet(object.InstID) ? String(object.InstID) : "",
      ID: isSet(object.ID) ? String(object.ID) : "",
      Type: isSet(object.Type) ? String(object.Type) : "",
    };
  },

  toJSON(message: EntityID): unknown {
    const obj: any = {};
    message.ProxyID !== undefined && (obj.ProxyID = message.ProxyID);
    message.InstID !== undefined && (obj.InstID = message.InstID);
    message.ID !== undefined && (obj.ID = message.ID);
    message.Type !== undefined && (obj.Type = message.Type);
    return obj;
  },

  create(base?: DeepPartial<EntityID>): EntityID {
    return EntityID.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EntityID>): EntityID {
    const message = createBaseEntityID();
    message.ProxyID = object.ProxyID ?? "";
    message.InstID = object.InstID ?? "";
    message.ID = object.ID ?? "";
    message.Type = object.Type ?? "";
    return message;
  },
};

function createBaseHeader(): Header {
  return {
    Version: 0,
    RequestId: "",
    Timeout: 0,
    RequestType: 0,
    ServiceType: 0,
    SrcEntity: undefined,
    DstEntity: undefined,
    Method: "",
    ContentType: 0,
    CompressType: 0,
    CheckFlags: 0,
    TransInfo: {},
    ErrCode: 0,
    ErrMsg: "",
  };
}

export const Header = {
  encode(message: Header, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Version !== 0) {
      writer.uint32(8).uint32(message.Version);
    }
    if (message.RequestId !== "") {
      writer.uint32(18).string(message.RequestId);
    }
    if (message.Timeout !== 0) {
      writer.uint32(24).uint32(message.Timeout);
    }
    if (message.RequestType !== 0) {
      writer.uint32(32).int32(message.RequestType);
    }
    if (message.ServiceType !== 0) {
      writer.uint32(40).int32(message.ServiceType);
    }
    if (message.SrcEntity !== undefined) {
      EntityID.encode(message.SrcEntity, writer.uint32(50).fork()).ldelim();
    }
    if (message.DstEntity !== undefined) {
      EntityID.encode(message.DstEntity, writer.uint32(58).fork()).ldelim();
    }
    if (message.Method !== "") {
      writer.uint32(66).string(message.Method);
    }
    if (message.ContentType !== 0) {
      writer.uint32(72).int32(message.ContentType);
    }
    if (message.CompressType !== 0) {
      writer.uint32(80).int32(message.CompressType);
    }
    if (message.CheckFlags !== 0) {
      writer.uint32(88).uint32(message.CheckFlags);
    }
    Object.entries(message.TransInfo).forEach(([key, value]) => {
      Header_TransInfoEntry.encode({ key: key as any, value }, writer.uint32(98).fork()).ldelim();
    });
    if (message.ErrCode !== 0) {
      writer.uint32(104).int32(message.ErrCode);
    }
    if (message.ErrMsg !== "") {
      writer.uint32(114).string(message.ErrMsg);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Header {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeader();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Version = reader.uint32();
          break;
        case 2:
          message.RequestId = reader.string();
          break;
        case 3:
          message.Timeout = reader.uint32();
          break;
        case 4:
          message.RequestType = reader.int32() as any;
          break;
        case 5:
          message.ServiceType = reader.int32() as any;
          break;
        case 6:
          message.SrcEntity = EntityID.decode(reader, reader.uint32());
          break;
        case 7:
          message.DstEntity = EntityID.decode(reader, reader.uint32());
          break;
        case 8:
          message.Method = reader.string();
          break;
        case 9:
          message.ContentType = reader.int32() as any;
          break;
        case 10:
          message.CompressType = reader.int32() as any;
          break;
        case 11:
          message.CheckFlags = reader.uint32();
          break;
        case 12:
          const entry12 = Header_TransInfoEntry.decode(reader, reader.uint32());
          if (entry12.value !== undefined) {
            message.TransInfo[entry12.key] = entry12.value;
          }
          break;
        case 13:
          message.ErrCode = reader.int32();
          break;
        case 14:
          message.ErrMsg = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Header {
    return {
      Version: isSet(object.Version) ? Number(object.Version) : 0,
      RequestId: isSet(object.RequestId) ? String(object.RequestId) : "",
      Timeout: isSet(object.Timeout) ? Number(object.Timeout) : 0,
      RequestType: isSet(object.RequestType) ? requestTypeFromJSON(object.RequestType) : 0,
      ServiceType: isSet(object.ServiceType) ? serviceTypeFromJSON(object.ServiceType) : 0,
      SrcEntity: isSet(object.SrcEntity) ? EntityID.fromJSON(object.SrcEntity) : undefined,
      DstEntity: isSet(object.DstEntity) ? EntityID.fromJSON(object.DstEntity) : undefined,
      Method: isSet(object.Method) ? String(object.Method) : "",
      ContentType: isSet(object.ContentType) ? contentTypeFromJSON(object.ContentType) : 0,
      CompressType: isSet(object.CompressType) ? compressTypeFromJSON(object.CompressType) : 0,
      CheckFlags: isSet(object.CheckFlags) ? Number(object.CheckFlags) : 0,
      TransInfo: isObject(object.TransInfo)
        ? Object.entries(object.TransInfo).reduce<{ [key: string]: Uint8Array }>((acc, [key, value]) => {
          acc[key] = bytesFromBase64(value as string);
          return acc;
        }, {})
        : {},
      ErrCode: isSet(object.ErrCode) ? Number(object.ErrCode) : 0,
      ErrMsg: isSet(object.ErrMsg) ? String(object.ErrMsg) : "",
    };
  },

  toJSON(message: Header): unknown {
    const obj: any = {};
    message.Version !== undefined && (obj.Version = Math.round(message.Version));
    message.RequestId !== undefined && (obj.RequestId = message.RequestId);
    message.Timeout !== undefined && (obj.Timeout = Math.round(message.Timeout));
    message.RequestType !== undefined && (obj.RequestType = requestTypeToJSON(message.RequestType));
    message.ServiceType !== undefined && (obj.ServiceType = serviceTypeToJSON(message.ServiceType));
    message.SrcEntity !== undefined &&
      (obj.SrcEntity = message.SrcEntity ? EntityID.toJSON(message.SrcEntity) : undefined);
    message.DstEntity !== undefined &&
      (obj.DstEntity = message.DstEntity ? EntityID.toJSON(message.DstEntity) : undefined);
    message.Method !== undefined && (obj.Method = message.Method);
    message.ContentType !== undefined && (obj.ContentType = contentTypeToJSON(message.ContentType));
    message.CompressType !== undefined && (obj.CompressType = compressTypeToJSON(message.CompressType));
    message.CheckFlags !== undefined && (obj.CheckFlags = Math.round(message.CheckFlags));
    obj.TransInfo = {};
    if (message.TransInfo) {
      Object.entries(message.TransInfo).forEach(([k, v]) => {
        obj.TransInfo[k] = base64FromBytes(v);
      });
    }
    message.ErrCode !== undefined && (obj.ErrCode = Math.round(message.ErrCode));
    message.ErrMsg !== undefined && (obj.ErrMsg = message.ErrMsg);
    return obj;
  },

  create(base?: DeepPartial<Header>): Header {
    return Header.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Header>): Header {
    const message = createBaseHeader();
    message.Version = object.Version ?? 0;
    message.RequestId = object.RequestId ?? "";
    message.Timeout = object.Timeout ?? 0;
    message.RequestType = object.RequestType ?? 0;
    message.ServiceType = object.ServiceType ?? 0;
    message.SrcEntity = (object.SrcEntity !== undefined && object.SrcEntity !== null)
      ? EntityID.fromPartial(object.SrcEntity)
      : undefined;
    message.DstEntity = (object.DstEntity !== undefined && object.DstEntity !== null)
      ? EntityID.fromPartial(object.DstEntity)
      : undefined;
    message.Method = object.Method ?? "";
    message.ContentType = object.ContentType ?? 0;
    message.CompressType = object.CompressType ?? 0;
    message.CheckFlags = object.CheckFlags ?? 0;
    message.TransInfo = Object.entries(object.TransInfo ?? {}).reduce<{ [key: string]: Uint8Array }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    message.ErrCode = object.ErrCode ?? 0;
    message.ErrMsg = object.ErrMsg ?? "";
    return message;
  },
};

function createBaseHeader_TransInfoEntry(): Header_TransInfoEntry {
  return { key: "", value: new Uint8Array() };
}

export const Header_TransInfoEntry = {
  encode(message: Header_TransInfoEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Header_TransInfoEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeader_TransInfoEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Header_TransInfoEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? bytesFromBase64(object.value) : new Uint8Array(),
    };
  },

  toJSON(message: Header_TransInfoEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : new Uint8Array()));
    return obj;
  },

  create(base?: DeepPartial<Header_TransInfoEntry>): Header_TransInfoEntry {
    return Header_TransInfoEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Header_TransInfoEntry>): Header_TransInfoEntry {
    const message = createBaseHeader_TransInfoEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? new Uint8Array();
    return message;
  },
};

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
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends { $case: string } ? { [K in keyof Omit<T, "$case">]?: DeepPartial<T[K]> } & { $case: T["$case"] }
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
