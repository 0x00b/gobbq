// type only
import type { UnaryRequestMessage, UnaryResponseMessage } from '../codec/msg';
import type { Endpoint } from '../endpoint';

interface Result {
  /** 远程接入点 */
  remote: Endpoint;
  /** 本地接入点 */
  local: Endpoint;
}

// export interface StreamResult extends Result {
//   /** Source 实例是用于操作流的句柄，流的后续操作通过 Source 完成 */
//   stream: Source;
//   /** Server 端返回的 Init Message */
//   init: StreamInitMessage;
// }

export interface UnaryResult extends Result {
  /** Server 端返回的 Message */
  response: UnaryResponseMessage;
}

/**
 * {@link ClientTransport} 绑定到特定的远程接入点，作为传输层负责收发 unary/stream message
 * - 对于 unary：实现 {@link addUnary } 接口
 * - 对于 stream：实现 {@link addStream } 接口
 *
 * 某些传输层可能不支持 stream 传输，但仍然要实现 {@link addStream} 接口，在方法内部抛出断言失败
 * 参考内置的 {@link UDPTransport} 实现
 */
export abstract class ClientTransport {
  /** 绑定到远程接入点 */
  protected abstract remote: Endpoint;

  /**
   * 收到服务端响应后 Promise 返回服务端响应，
   * 以及本次传输使用的本地接入点和远程接入点
   */
  public abstract addUnary(req: UnaryRequestMessage): Promise<UnaryResult> | UnaryResult;
  /**
   * 流的 Init 阶段完成后 Promise 返回服务端 Init 响应、
   * 用于后续流操作的句柄，
   * 以及本次传输使用的本地接入点和远程接入点
   */
  // public abstract addStream(req: StreamInitMessage): Promise<StreamResult> | StreamResult;

  /**
   * 主动销毁
   */
  public abstract destroy(): void;

  /**
   * 连接到远程接入点
   * @remark
   * 某些传输层可能不需要”连接“，但仍然要实现此接口，返回 `Promise.promise()` 即可，
   * 参考内置的 实现
   */
   public abstract connect(): Promise<void>;


   public abstract send(buffer: Buffer):void;

}

