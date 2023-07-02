// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

import { UnaryResponse } from "gobbq-ts/dist/src/context/unary";
import { Client, Context } from "gobbq-ts/dist/src";
import { makeClientConstructor } from "gobbq-ts/dist/src/bbq/bbq";
import { ServiceType } from "gobbq-ts/dist/proto/bbq";
import Long from "long";
import { SayHelloRequest } from "./exam"
import { SayHelloResponse } from "./exam"
	
// EchoService
export type EchoServiceDefinition = typeof EchoServiceDefinition;
export const EchoServiceDefinition = {
  typeName: "exampb.EchoService",
  serviceType: ServiceType.Service, 
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: SayHelloResponse,
      requestSerialize: (req: SayHelloRequest): Buffer => {
        return Buffer.from(SayHelloRequest.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): SayHelloRequest => {
        return SayHelloRequest.decode(input)
      },
      responseSerialize: (req: SayHelloResponse): Buffer => {
        return Buffer.from(SayHelloResponse.encode(req).finish())
      },
      responseDeserialize: (input: Uint8Array): SayHelloResponse => {
        return SayHelloResponse.decode(input)
      },
    },
	},
} as const;

export interface EchoService {

	// SayHello
	SayHello(request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export interface EchoServiceService {

	// SayHello
	SayHello(c: Context, request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export function RegisterEchoServiceService(client: Client<any>, svc: EchoServiceService){
  client.register(EchoServiceDefinition, svc)
}

export function NewEchoService(client: Client<any>): EchoService {
  return makeClientConstructor(client, EchoServiceDefinition) as unknown as EchoService
}
// EchoEtyEntity
export type EchoEtyEntityDefinition = typeof EchoEtyEntityDefinition;
export const EchoEtyEntityDefinition = {
  typeName: "exampb.EchoEtyEntity",
  serviceType: ServiceType.Entity, 
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: SayHelloResponse,
      requestSerialize: (req: SayHelloRequest): Buffer => {
        return Buffer.from(SayHelloRequest.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): SayHelloRequest => {
        return SayHelloRequest.decode(input)
      },
      responseSerialize: (req: SayHelloResponse): Buffer => {
        return Buffer.from(SayHelloResponse.encode(req).finish())
      },
      responseDeserialize: (input: Uint8Array): SayHelloResponse => {
        return SayHelloResponse.decode(input)
      },
    },
	},
} as const;

export interface EchoEtyEntity {

	// SayHello
	SayHello(request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export interface EchoEtyEntityService {

	// SayHello
	SayHello(c: Context, request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export function RegisterEchoEtyEntityService(client: Client<any>, svc: EchoEtyEntityService){
  client.register(EchoEtyEntityDefinition, svc)
}

export function NewEchoEtyEntity(client: Client<any>, entityID: Long): EchoEtyEntity {
  return makeClientConstructor(client, EchoEtyEntityDefinition, entityID) as unknown as EchoEtyEntity
}

	
// EchoSvc2Service
export type EchoSvc2ServiceDefinition = typeof EchoSvc2ServiceDefinition;
export const EchoSvc2ServiceDefinition = {
  typeName: "exampb.EchoSvc2Service",
  serviceType: ServiceType.Service, 
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: SayHelloResponse,
      requestSerialize: (req: SayHelloRequest): Buffer => {
        return Buffer.from(SayHelloRequest.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): SayHelloRequest => {
        return SayHelloRequest.decode(input)
      },
      responseSerialize: (req: SayHelloResponse): Buffer => {
        return Buffer.from(SayHelloResponse.encode(req).finish())
      },
      responseDeserialize: (input: Uint8Array): SayHelloResponse => {
        return SayHelloResponse.decode(input)
      },
    },
	},
} as const;

export interface EchoSvc2Service {

	// SayHello
	SayHello(request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export interface EchoSvc2ServiceService {

	// SayHello
	SayHello(c: Context, request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export function RegisterEchoSvc2ServiceService(client: Client<any>, svc: EchoSvc2ServiceService){
  client.register(EchoSvc2ServiceDefinition, svc)
}

export function NewEchoSvc2Service(client: Client<any>): EchoSvc2Service {
  return makeClientConstructor(client, EchoSvc2ServiceDefinition) as unknown as EchoSvc2Service
}
// ClientEntity 客户端
export type ClientEntityDefinition = typeof ClientEntityDefinition;
export const ClientEntityDefinition = {
  typeName: "exampb.ClientEntity",
  serviceType: ServiceType.Entity, 
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: SayHelloResponse,
      requestSerialize: (req: SayHelloRequest): Buffer => {
        return Buffer.from(SayHelloRequest.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): SayHelloRequest => {
        return SayHelloRequest.decode(input)
      },
      responseSerialize: (req: SayHelloResponse): Buffer => {
        return Buffer.from(SayHelloResponse.encode(req).finish())
      },
      responseDeserialize: (input: Uint8Array): SayHelloResponse => {
        return SayHelloResponse.decode(input)
      },
    },
	},
} as const;

export interface ClientEntity {

	// SayHello
	SayHello(request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export interface ClientEntityService {

	// SayHello
	SayHello(c: Context, request: SayHelloRequest):UnaryResponse<SayHelloResponse>
}

export function RegisterClientEntityService(client: Client<any>, svc: ClientEntityService){
  client.register(ClientEntityDefinition, svc)
}

export function NewClientEntity(client: Client<any>, entityID: Long): ClientEntity {
  return makeClientConstructor(client, ClientEntityDefinition, entityID) as unknown as ClientEntity
}
// NoRespEntity 客户端
export type NoRespEntityDefinition = typeof NoRespEntityDefinition;
export const NoRespEntityDefinition = {
  typeName: "exampb.NoRespEntity",
  serviceType: ServiceType.Entity, 
  methods: {
    SayHello: {
      methodName: "SayHello",
      requestType: SayHelloRequest,
      responseType: undefined,
      requestSerialize: (req: SayHelloRequest): Buffer => {
        return Buffer.from(SayHelloRequest.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): SayHelloRequest => {
        return SayHelloRequest.decode(input)
      },
      responseSerialize: (req: any): Buffer => {
        return Buffer.from("")
      },
      responseDeserialize: (input: Uint8Array): any => {
        
      },
    },
	},
} as const;

export interface NoRespEntity {

	// SayHello
	SayHello(request: SayHelloRequest):void
}

export interface NoRespEntityService {

	// SayHello
	SayHello(c: Context, request: SayHelloRequest):void
}

export function RegisterNoRespEntityService(client: Client<any>, svc: NoRespEntityService){
  client.register(NoRespEntityDefinition, svc)
}

export function NewNoRespEntity(client: Client<any>, entityID: Long): NoRespEntity {
  return makeClientConstructor(client, NoRespEntityDefinition, entityID) as unknown as NoRespEntity
}

