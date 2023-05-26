// GENERATED CODE -- DO NOT EDIT!

// package: exampb
// file: exam.proto

import * as exam_pb from "./exam_pb";
import * as grpc from "grpc";

interface IEchoService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  sayHello: grpc.MethodDefinition<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export const EchoService: IEchoService;

export interface IEchoServer extends grpc.UntypedServiceImplementation {
  sayHello: grpc.handleUnaryCall<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export class EchoClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  sayHello(argument: exam_pb.SayHelloRequest, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
}

interface IEchoEtyService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  sayHello: grpc.MethodDefinition<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export const EchoEtyService: IEchoEtyService;

export interface IEchoEtyServer extends grpc.UntypedServiceImplementation {
  sayHello: grpc.handleUnaryCall<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export class EchoEtyClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  sayHello(argument: exam_pb.SayHelloRequest, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
}

interface IEchoSvc2Service extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  sayHello: grpc.MethodDefinition<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export const EchoSvc2Service: IEchoSvc2Service;

export interface IEchoSvc2Server extends grpc.UntypedServiceImplementation {
  sayHello: grpc.handleUnaryCall<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export class EchoSvc2Client extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  sayHello(argument: exam_pb.SayHelloRequest, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
}

interface IClientService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  sayHello: grpc.MethodDefinition<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export const ClientService: IClientService;

export interface IClientServer extends grpc.UntypedServiceImplementation {
  sayHello: grpc.handleUnaryCall<exam_pb.SayHelloRequest, exam_pb.SayHelloResponse>;
}

export class ClientClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  sayHello(argument: exam_pb.SayHelloRequest, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
  sayHello(argument: exam_pb.SayHelloRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<exam_pb.SayHelloResponse>): grpc.ClientUnaryCall;
}
