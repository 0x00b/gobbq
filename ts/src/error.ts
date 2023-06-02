/**
 * RPC 错误码
 *
 * @remark
 * 使用 {@link https://www.typescriptlang.org/docs/handbook/enums.html#const-enums | const enum} 实现常量，编译时内联
 */
export const enum ERROR {
  
  CLIENT_INVALID_ERR = 1,
  /** 客户端调用超时 */
  CLIENT_INVOKE_TIMEOUT_ERR = 2,

  /** 客户端连接错误 */
  CLIENT_CONNECT_ERR = 3,

  /** 客户端编码错误 */
  CLIENT_ENCODE_ERR = 4,
  /** 客户端解码错误 */
  CLIENT_DECODE_ERR = 5,

  /** 客户端选ip路由错误 */
  CLIENT_ROUTER_ERR = 6,

  /** 客户端网络错误 */
  CLIENT_NETWORK_ERR = 7,

  /** 客户端响应参数自动校验失败错误 */
  CLIENT_VALIDATE_ERR = 8,

  /** 上游主动断开连接，提前取消请求错误 */
  CLIENT_CANCELED_ERR = 9,

  /** 客户端流式网络错误, 详细错误码需要在实现过程中再梳理 */
  STREAM_CLIENT_NETWORK_ERR = 10,

  /** 客户端流式传输错误, 详细错误码需要在实现过程中再梳理；比如：流消息过大等 */
  STREAM_CLIENT_MSG_EXCEED_LIMIT_ERR = 11,

  /** 客户端流式编码错误 */
  STREAM_CLIENT_ENCODE_ERR = 12,
  /** 客户端流式编解码错误 */
  STREAM_CLIENT_DECODE_ERR = 13,

  // 客户端流式写错误, 详细错误码需要在实现过程中再梳理
  STREAM_CLIENT_WRITE_END = 331,
  STREAM_CLIENT_WRITE_OVERFLOW_ERR = 332,
  STREAM_CLIENT_WRITE_CLOSE_ERR = 333,
  STREAM_CLIENT_WRITE_TIMEOUT_ERR = 334,

  // 客户端流式读错误, 详细错误码需要在实现过程中再梳理
  STREAM_CLIENT_READ_END = 351,
  STREAM_CLIENT_READ_CLOSE_ERR = 352,
  STREAM_CLIENT_READ_EMPTY_ERR = 353,
  STREAM_CLIENT_READ_TIMEOUT_ERR = 354,

  SERVER_TIMEOUT_ERR = 400,
  
  /** 未明确的错误 */
  INVOKE_UNKNOWN_ERR = 999,
  /** 未明确的错误 */
  STREAM_UNKNOWN_ERR = 1000,
}

export class RpcError implements Error {
  public readonly name: string;
  public readonly stack?: string;

  public constructor(
    public readonly code: ERROR ,
    public readonly message: string,
    public readonly cause?: any,
  ) {
    this.name = `[ Error ${code} ${code}] ${message}`;

    const error = new Error(this.name);
    Error.captureStackTrace(error, RpcError);
    this.stack = error.stack;
  }

  public toJSON() {
    return {
      pid: process.pid,
      code: this.code,
      message: this.message,
      name: this.name,
      stack: this.stack,
    };
  }
};
