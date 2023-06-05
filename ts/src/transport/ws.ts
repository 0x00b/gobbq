import * as assert from 'assert';

import { Deferred, noop } from '../utils';
import { ERROR, RpcError } from '../error';
import { ClientTransport, UnaryResult } from './base';
import { WebSocket } from 'ws';

import {
  decode,
} from '../codec/msg';

// import type { Transport } from '@/-stream';
// import type { StreamResult, UnaryResult } from './base';
import type { Endpoint } from '../endpoint';
import { createTruncator } from './truncator';
import { Packet } from '../codec/packet';

// eslint-disable-next-line @typescript-eslint/no-require-imports
const DEBUGGER = require('debug')(':client:ws');

/**
 * 内置的 WS  传输实现
 *
 * - 使用 WS 传输 unary 时，需要根据 requestId 来匹配 response 和 request
 * - 使用 WS 传输 stream 时，需要根据 streamId 来区分不通 stream 上帧
 *
 * {@link WSTransport} 实现了 {@link Transport} 的 send 接口，使用 Source 管理流
 * - 向 socket 中发送帧时：Source -> Transport：send(buffer)
 * - 从 socket 中切分出帧时：Transport -> Source: push(buffer)
 */
export class WSTransport extends ClientTransport /*implements Transport*/ {
  /** INIT 完成的 stream */
  // private streams = new Map<number, Source>();
  // /** INIT 阶段的 stream */
  // private deferredStream = new Map<number, Deferred<StreamResult, RpcError>>();
  /** 未收到回包的 unary  */
  private deferredUnary = new Map<string, Deferred<UnaryResult, RpcError>>();

  private socket?: WebSocket;
  private localpoint: Endpoint = {
    host: 'NOHOST',
    port: 0,
    protocol: 'ws',
  };

  private promise?: Promise<void>;
  private connected = false;
  private destroyed = false;
  public constructor(
    protected remotepoint: Endpoint,
    private readonly onDestroyed: () => void,
    private readonly onUnaryMessage: (pkt: Packet) => void,
  ) {
    super();
  }

  /**
   * 主动销毁
   */
  public destroy(error?: RpcError) {
    DEBUGGER('[destroy]', `destroyed:${this.destroyed}`);
    /* istanbul ignore if */
    if (this.destroyed) return;
    this.destroyed = true;
    this.connected = false;

    // 清理请求
    this.clear(error ?? new RpcError(ERROR.CLIENT_CANCELED_ERR, 'Destroyed'));

    // 清理 socket
    this.socket?.close();
    this.socket = undefined;
    this.promise = undefined;

    this.onDestroyed();
  }

  public connect() {
    DEBUGGER('[connect]', `connected:${this.connected}`);

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
      const { port, host, protocol } = this.remotepoint;
      var serverAddr = `${protocol}://${host}:${port}`
      var httporigin = `http://${host}:${port}/`
      console.log(serverAddr)
      const websocket = new WebSocket(serverAddr, { origin: httporigin })
      websocket.binaryType = 'arraybuffer'

      websocket.onerror = function (e) {
        console.log(e)
        reject(new RpcError(ERROR.CLIENT_CONNECT_ERR, "WebSocket连接发生错误"));
        return
      }

      let self = this

      //连接成功建立的回调方法
      let open = function () {
        console.log("WebSocket连接成功")
        const handleData = createTruncator(self.onFrame.bind(self));
        //接收到消息的回调方法
        websocket.onmessage = function (event) {
          var data = event.data
          const buffer = new Uint8Array(data as any)
          // console.log("收到数据：",  data)
          // handleData(data)
          self.onData(websocket, handleData, Buffer.from(buffer))
        }

        //连接关闭的回调方法
        websocket.onclose = function () {
          console.log("WebSocket连接关闭")
          self.onClose(true)
        }
        self.connected = true
        resolve();
      }
      websocket.onopen = open

      this.socket = websocket;
    });

    return this.promise;
  }

  public send(buffer: Buffer) {
    DEBUGGER('[send]', `connected:${this.connected}`, `promise:${this.promise}`);

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

    return this.socket.send(buffer);
  }

  private onData(socket: WebSocket, handleData: (chunk: Buffer) => void, buffer: Buffer) {
    try {
      handleData(buffer);
    } catch (error) {
      // 触发 'error' 和 'close'
      socket.close();
    }
  }

  /**
   * 不论是主动 destroy 还是被动关闭 socket，最终都会触发 `close`
   * 此时 `socket.destroyed === true`
   */
  private onClose(hadError: boolean) {
    DEBUGGER('[close]', `connected:${this.connected}`, `hasError:${hadError}`, `destroyed:${this.socket?.close}`, `local:${this.local}`);
    this.connected = false;
    this.socket = undefined;
    this.promise = undefined;

    // 清理请求
    this.clear(new RpcError(ERROR.CLIENT_CANCELED_ERR, 'Closed'));
  }

  /**
   * socket 发生错误
   * - 解码错误
   * - ECONNREFUSED 连接不到对端
   * 之后会触发 'close'
   * @param err Error
   */
  private onError(err: Error) {
    DEBUGGER('[error]', err);
    this.connected = false;
    this.clear(err instanceof RpcError ? err : new RpcError(ERROR.CLIENT_CONNECT_ERR, err.message, err));
    // 清理 socket
    this.socket?.close();
    this.socket = undefined;
    this.promise = undefined;
  }

  private clear(err: RpcError) {
    // INIT 完成的 stream
    // this.streams.forEach((source) => {
    //   source.destroy(err);
    // });

    // this.streams.clear();

    // INIT 阶段的 stream
    // if (this.deferredStream.size > 0) {
    //   this.deferredStream.forEach((d) => {
    //     d.reject(err);
    //   });
    //   this.deferredStream.clear();
    // }

    // 未收到回包的 unary
    if (this.deferredUnary.size > 0) {
      this.deferredUnary.forEach((d) => {
        d.reject(err);
      });
      this.deferredUnary.clear();
    }
  }

  /**
   * 处理从字节流中切分出的帧
   * @param buffer 帧
   */
  private onFrame(buffer: Buffer) {
    const pkt = decode(buffer);

    // if (isUnaryMessage(message)) {
    DEBUGGER('[receive] unary', `requestId:${pkt?.Header.RequestId}`);
    this.onUnaryMessage(pkt as Packet);
    return;
    // }

    // if (isStreamMessage(message)) {
    //   DEBUGGER('[receive] stream', `streamId:${message.streamId} streamType:${message.streamType}`);
    //   this.onStreamMessage(message as StreamMessage);
    // }
  }

  public local(): Endpoint {
    return this.localpoint
  }

  public remote(): Endpoint {
    return this.remotepoint
  }

}
