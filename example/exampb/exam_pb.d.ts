// package: exampb
// file: exam.proto

import * as jspb from "google-protobuf";
import * as bbq_pb from "./bbq_pb";

export class SayHelloRequest extends jspb.Message {
  getText(): string;
  setText(value: string): void;

  hasClientid(): boolean;
  clearClientid(): void;
  getClientid(): bbq_pb.EntityID | undefined;
  setClientid(value?: bbq_pb.EntityID): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SayHelloRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SayHelloRequest): SayHelloRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SayHelloRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SayHelloRequest;
  static deserializeBinaryFromReader(message: SayHelloRequest, reader: jspb.BinaryReader): SayHelloRequest;
}

export namespace SayHelloRequest {
  export type AsObject = {
    text: string,
    clientid?: bbq_pb.EntityID.AsObject,
  }
}

export class SayHelloResponse extends jspb.Message {
  getText(): string;
  setText(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SayHelloResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SayHelloResponse): SayHelloResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SayHelloResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SayHelloResponse;
  static deserializeBinaryFromReader(message: SayHelloResponse, reader: jspb.BinaryReader): SayHelloResponse;
}

export namespace SayHelloResponse {
  export type AsObject = {
    text: string,
  }
}

