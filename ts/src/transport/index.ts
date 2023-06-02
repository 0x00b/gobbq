/**
 *  Node.js 客户端传输层
 *
 * {@link ClientTransport} 绑定到特定的远程接入点，作为传输层负责收发 unary/stream message
 *
 * {@link ClientTransport} 是一个 `abstract class`，内置了三个实现：
 * > {@link TCPTransport} | {@link UDPTransport} | {@link HTTPTransport}
 *
 * @packageDocumentation
 */
// type only
import { Kcp } from 'kcpjs';
import { Packet } from '../codec/packet';
import type { Endpoint } from '../endpoint';
import type { /*StreamResult,*/ UnaryResult } from './base';

// builtin implements
import { ClientTransport } from './base';
import { KCPTransport } from './kcp';
import { TCPTransport } from './tcp';
import { WSTransport } from './ws';
// import { WSTransport } from './ws';

export type CreateTransport = typeof createTransport;
export type { /*StreamResult,*/ UnaryResult };
export { TCPTransport, ClientTransport };

/**
 * 使用内置 {@link ClientTransport } 实现
 * @param endPoint 远程接入点，根据 {@link Endpoint.protocol } 自动选择具体的传输实现
 * @param onDestroyed 销毁时的回调
 */
export function createTransport(endPoint: Endpoint, onDestroyed: () => void, onUnaryMessage: (pkt:Packet) => void,): ClientTransport {
  switch (endPoint.protocol) {
    case 'kcp':
      return new KCPTransport(endPoint, onDestroyed, onUnaryMessage);
    case 'ws':
    case 'wss':
    return new WSTransport(endPoint, onDestroyed, onUnaryMessage);
    default: // 默认使用 tcp
    case 'tcp':
      return new TCPTransport(endPoint, onDestroyed, onUnaryMessage);
  }
};
