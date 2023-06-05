import * as assert from 'assert';

import { RpcError } from '../error';
import { ClientTransport } from './base';

import {
  decode,
} from '../codec/msg';

import type { Endpoint } from '../endpoint';
import { createTruncator } from './truncator';
import { Packet } from '../codec/packet';
import { KCP } from './node-kcp-ex/kcpclient';

const block = undefined

/**
 * 内置的 KCP  传输实现
 *
 * - 使用 KCP 传输 unary 时，需要根据 requestId 来匹配 response 和 request
 * - 使用 KCP 传输 stream 时，需要根据 streamId 来区分不通 stream 上帧
 *
 * {@link KCPTransport} 实现了 {@link Transport} 的 send 接口，使用 Source 管理流
 * - 向 socket 中发送帧时：Source -> Transport：send(buffer)
 * - 从 socket 中切分出帧时：Transport -> Source: push(buffer)
 */
export class KCPTransport extends ClientTransport /*implements Transport*/ {

  private socket?: KCP;
  private localpoint: Endpoint = {
    host: 'NOHOST',
    port: 0,
    protocol: 'kcp',
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
    console.log('[destroy]', `destroyed:${this.destroyed}`);
    /* istanbul ignore if */
    if (this.destroyed) return;
    this.destroyed = true;
    this.connected = false;

    // 清理请求
    // this.clear(error ?? new RpcError(ERROR.CLIENT_CANCELED_ERR, 'Destroyed'));

    // 清理 socket
    this.socket?.close();
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

      const handleData = createTruncator(this.onFrame.bind(this));

      const socket = new KCP(255, {
        port: port,
        address: host,
        onrecv: this.onData.bind(this, handleData),
        onerror: this.onError.bind(this),
      })

      resolve()
      this.connected = true;

      if (!socket) {
        reject("dial failed")
        return
      }
      // socket.setACKNoDelay(true)
      // socket.setStreamMode(true)
      // socket.setNoDelay(1,10,2,1)
      // socket.check
      // socket
      //   .on('error', this.onError.bind(this))
      //   .on('close', this.onClose.bind(this))
      //   // .on('end', this.onEnd.bind(this))
      //   .on('recv', this.onData.bind(this, socket, handleData))
      //   // .on("connection",()=>{
      //     console.log("connnnn")
      //     resolve()
      //     this.connected = true;
      // })

      // if (socket.localAddress && socket.localPort ) {
      //   this.local = {
      //     host: socket.localAddress,
      //     port: socket.localPort,
      //     protocol: 'kcp',
      //   };
      // }
      // resolve();

      this.socket = socket;
    })

    return this.promise;
  }

  private onData(handleData: (chunk: Buffer) => void, buffer: Buffer,/* len:any*/) {
    try {
      console.log("recv buffer", buffer)
      handleData(buffer);
    } catch (error) {
      // 触发 'error' 和 'close'
      console.log("ondata err:", error)
      this.socket?.close();
    }
  }

  /**
   * 处理从字节流中切分出的帧
   * @param buffer 帧
   */
  private onFrame(buffer: Buffer) {
    // console.log("recv:", buffer)
    const pkt = decode(buffer);

    // if (isUnaryMessage(message)) {
    console.log('[receive] unary', `requestId:${pkt?.Header.RequestId}`);
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

    return this.socket.send(buffer);
  }

  /**
   * 不论是主动 destroy 还是被动关闭 socket，最终都会触发 `close`
   * 此时 `socket.destroyed === true`
   */
  private onClose(hadError: boolean) {
    console.log('[close]', `connected:${this.connected}`, `hasError:${hadError}`, `destroyed:${this.socket?.close}`, `local:${this.local}`);
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
    this.socket?.close();
    this.socket = undefined;
    this.promise = undefined;
  }


  public local(): Endpoint {
    return this.localpoint
  }

  public remote(): Endpoint {
    return this.remotepoint
  }

}
