import * as net from 'net';
import * as assert from 'assert';

import { ERROR, RpcError } from '../error';
import { ClientTransport } from './base';


import {
  decode,

} from '../codec/msg';

// import type { Transport } from '@/-stream';
// import type { StreamResult, UnaryResult } from './base';
import type { Endpoint } from '../endpoint';
import { createTruncator } from './truncator';
import { Packet } from '../codec/packet';

// eslint-disable-next-line @typescript-eslint/no-require-imports
// const console.log = require('debug')(':client:tcp');

/**
 * 内置的 TCP  传输实现
 *
 * - 使用 TCP 传输 unary 时，需要根据 requestId 来匹配 response 和 request
 * - 使用 TCP 传输 stream 时，需要根据 streamId 来区分不通 stream 上帧
 *
 * {@link TCPTransport} 实现了 {@link Transport } 的 send 接口，使用 Source 管理流
 * - 向 socket 中发送帧时：Source -> Transport：send(buffer)
 * - 从 socket 中切分出帧时：Transport -> Source: push(buffer)
 */
export class TCPTransport extends ClientTransport /*implements Transport*/ {
  /** INIT 完成的 stream */
  // private streams = new Map<number, Source>();
  // /** INIT 阶段的 stream */
  // private deferredStream = new Map<number, Deferred<StreamResult, RpcError>>();

  private socket?: net.Socket;
  private localpoint: Endpoint = {
    host: 'NOHOST',
    port: 0,
    protocol: 'tcp',
  };

  private promise?: Promise<void>;
  private connected = false;
  private destroyed = false;
  public constructor(
    protected remotepoint: Endpoint,
    private readonly onDestroyed: () => void,
    private readonly onUnaryMessage: (pkt:Packet) => void,
  ) {
    super();
  }

  /**
   * 主动销毁
   */
  public destroy(error?: RpcError) {
    console.log('[destroy]', `destroyed:${this.destroyed}`);
    /* istanbul ignore if */
    if (this.destroyed) return;
    this.destroyed = true;
    this.connected = false;

    // 清理请求
    // this.clear(error ?? new RpcError(ERROR.CLIENT_CANCELED_ERR, 'Destroyed'));

    // 清理 socket
    this.socket?.destroy();
    this.socket = undefined;
    this.promise = undefined;

    this.onDestroyed();
  }

  public connect() {
    console.log('[connect]', `connected:${this.connected}`);

    // 正在连接或已连接
    if (this.promise !== undefined) {
      return this.promise;
    }

    /* istanbul ignore if */
    if (this.socket !== undefined) {
      throw new assert.AssertionError({
        expected: undefined,
        actual: this.socket,
      });
    }

    this.promise = new Promise((resolve, reject) => {
      const { port, host } = this.remotepoint;
      const socket = net.createConnection({ port, host })
        .once('error', (error) => {
          reject(new RpcError(ERROR.CLIENT_NETWORK_ERR, error.message, error));
        })
        .setNoDelay(true)
        .once('connect', () => {
          socket.removeListener('error', reject);

          const handleData = createTruncator(this.onFrame.bind(this));
          socket
            .on('error', this.onError.bind(this))
            .on('close', this.onClose.bind(this))
            // .on('end', this.onEnd.bind(this))
            .on('data', this.onData.bind(this, socket, handleData));

          this.connected = true;

          if (socket.localAddress && socket.localPort) {
            this.localpoint = {
              host: socket.localAddress,
              port: socket.localPort,
              protocol: 'tcp',
            };
          }

          resolve();
        });

      this.socket = socket;
    });

    return this.promise;
  }

  private onData(socket: net.Socket, handleData: (chunk: Buffer) => void, buffer: Buffer) {
    try {
      handleData(buffer);
    } catch (error) {
      // 触发 'error' 和 'close'
      if (error instanceof Error) {
        socket.destroy(new RpcError(ERROR.CLIENT_DECODE_ERR, error.message, error));
      } else {
        socket.destroy(new RpcError(ERROR.CLIENT_DECODE_ERR, JSON.stringify(error)));
      }
    }
  }

  /**
   * 不论是主动 destroy 还是被动关闭 socket，最终都会触发 `close`
   * 此时 `socket.destroyed === true`
   */
  private onClose(hadError: boolean) {
    console.log('[close]', `connected:${this.connected}`, `hasError:${hadError}`, `destroyed:${this.socket?.destroyed}`, `local:${this.local}`);
    this.connected = false;
    this.socket = undefined;
    this.promise = undefined;

    // 清理请求
    // this.clear(new RpcError(ERROR.CLIENT_CANCELED_ERR, 'Closed'));
  }

  /**
   * socket 发生错误
   * - 解码错误
   * - ECONNREFUSED 连接不到对端
   * 之后会触发 'close'
   * @param err Error
   */
  private onError(err: Error) {
    console.log('[error]', err);
    this.connected = false;
    // this.clear(err instanceof RpcError ? err : new RpcError(ERROR.CLIENT_CONNECT_ERR, err.message, err));
    // 清理 socket
    this.socket?.destroy();
    this.socket = undefined;
    this.promise = undefined;
  }


  /**
   * 处理从字节流中切分出的帧
   * @param buffer 帧
   */
  private onFrame(buffer: Buffer) {
    // console.log("recv:", buffer)
    const pkt = decode(buffer);

    // if (isUnaryMessage(message)) {
    // console.log('[receive] unary', `requestId:${pkt?.Header.RequestId}`);
    this.onUnaryMessage(pkt as Packet);
    return;
    // }

    // if (isStreamMessage(message)) {
    //   console.log('[receive] stream', `streamId:${message.streamId} streamType:${message.streamType}`);
    //   this.onStreamMessage(message as StreamMessage);
    // }
  }
  
  public send(buffer: Buffer) {
    console.log('[send]', `connected:${this.connected}`, `promise:${this.promise}`);

    if (!this.connected) {
      return false;
    }

    /* istanbul ignore if */
    if (this.socket === undefined) {
      throw new assert.AssertionError({
        expected: '<net.Socket>',
        actual: undefined,
      });
    }

    return this.socket.write(buffer);
  }
  
  public  local(): Endpoint{
    return this.localpoint
  }

  public  remote(): Endpoint{
    return this.remotepoint
  }

}
