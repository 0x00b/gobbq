import { SayHelloRequest } from "../../../example/exampb/exam";
import { ServiceType } from "../../../proto/bbq/bbq";
import { Context } from "./context";

/**
 * A deserialization function
 * @param data The byte sequence to deserialize
 * @return The data deserialized as a value
 */
type deserialize<T> = (data: Buffer) => T;

/**
 * A serialization function
 * @param value The value to serialize
 * @return The value serialized as a byte sequence
 */
type serialize<T> = (value: T) => Buffer;

/**
 * An object that completely defines a service method signature.
 */
export interface MethodDefinition<RequestType, ResponseType> {
  /**
   * The method's URL path
   */
  methodName: string;

  requestType: RequestType,
  responseType: ResponseType,

  /**
  * Serialization function for request values
  */
  requestSerialize: serialize<RequestType>;
  /**
   * Serialization function for response values
   */
  responseSerialize: serialize<ResponseType>;
  /**
   * Deserialization function for request data
   */
  requestDeserialize: deserialize<RequestType>;
  /**
   * Deserialization function for repsonse data
   */
  responseDeserialize: deserialize<ResponseType>;

}

export interface MethodImpl<RequestType, ResponseType> extends MethodDefinition<RequestType, ResponseType> {
  handle: (ctx: Context, req: RequestType) => ResponseType | Error;
}

/**
 * An object that completely defines a service.
 */
export interface ServiceDefinition {
  typeName: string;
  serviceType: ServiceType;
  methods: Record<string, MethodDefinition<any, any>>;
}
