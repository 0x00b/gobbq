// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var exam_pb = require('./exam_pb.js');
var bbq_pb = require('./bbq_pb.js');

function serialize_exampb_SayHelloRequest(arg) {
  if (!(arg instanceof exam_pb.SayHelloRequest)) {
    throw new Error('Expected argument of type exampb.SayHelloRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_exampb_SayHelloRequest(buffer_arg) {
  return exam_pb.SayHelloRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_exampb_SayHelloResponse(arg) {
  if (!(arg instanceof exam_pb.SayHelloResponse)) {
    throw new Error('Expected argument of type exampb.SayHelloResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_exampb_SayHelloResponse(buffer_arg) {
  return exam_pb.SayHelloResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var EchoService = exports.EchoService = {
  sayHello: {
    path: '/exampb.Echo/SayHello',
    requestStream: false,
    responseStream: false,
    requestType: exam_pb.SayHelloRequest,
    responseType: exam_pb.SayHelloResponse,
    requestSerialize: serialize_exampb_SayHelloRequest,
    requestDeserialize: deserialize_exampb_SayHelloRequest,
    responseSerialize: serialize_exampb_SayHelloResponse,
    responseDeserialize: deserialize_exampb_SayHelloResponse,
  },
};

exports.EchoClient = grpc.makeGenericClientConstructor(EchoService);
var EchoEtyService = exports.EchoEtyService = {
  sayHello: {
    path: '/exampb.EchoEty/SayHello',
    requestStream: false,
    responseStream: false,
    requestType: exam_pb.SayHelloRequest,
    responseType: exam_pb.SayHelloResponse,
    requestSerialize: serialize_exampb_SayHelloRequest,
    requestDeserialize: deserialize_exampb_SayHelloRequest,
    responseSerialize: serialize_exampb_SayHelloResponse,
    responseDeserialize: deserialize_exampb_SayHelloResponse,
  },
};

exports.EchoEtyClient = grpc.makeGenericClientConstructor(EchoEtyService);
var EchoSvc2Service = exports.EchoSvc2Service = {
  sayHello: {
    path: '/exampb.EchoSvc2/SayHello',
    requestStream: false,
    responseStream: false,
    requestType: exam_pb.SayHelloRequest,
    responseType: exam_pb.SayHelloResponse,
    requestSerialize: serialize_exampb_SayHelloRequest,
    requestDeserialize: deserialize_exampb_SayHelloRequest,
    responseSerialize: serialize_exampb_SayHelloResponse,
    responseDeserialize: deserialize_exampb_SayHelloResponse,
  },
};

exports.EchoSvc2Client = grpc.makeGenericClientConstructor(EchoSvc2Service);
// 客户端
var ClientService = exports.ClientService = {
  sayHello: {
    path: '/exampb.Client/SayHello',
    requestStream: false,
    responseStream: false,
    requestType: exam_pb.SayHelloRequest,
    responseType: exam_pb.SayHelloResponse,
    requestSerialize: serialize_exampb_SayHelloRequest,
    requestDeserialize: deserialize_exampb_SayHelloRequest,
    responseSerialize: serialize_exampb_SayHelloResponse,
    responseDeserialize: deserialize_exampb_SayHelloResponse,
  },
};

exports.ClientClient = grpc.makeGenericClientConstructor(ClientService);
