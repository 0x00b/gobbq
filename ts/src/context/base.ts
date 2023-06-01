import { Deferred } from '../utils';
import type { RpcError } from '../error';
// type only
import type {
  // ,
  // StreamInitMessage,
  UnaryRequestMessage,
  UnaryResponseMessage,
  UnaryResponsePacket,
} from '../codec/msg';
import type { Endpoint } from '../endpoint';
import { CompressType, ContentType } from '../../../proto/bbq/bbq';
import { Packet } from '../codec/packet';
import { Middleware } from '../middleware';
import * as ds from '../dispatcher/compose';

// eslint-disable-next-line @typescript-eslint/no-require-imports
export const DEBUGGER = require('debug')(':client:ctx');

export interface Options {
  /**
   * 请求超时时间
   */
  timeout: number;

  /**
   * 框架信息透传的消息类型
   */
  messageType: number;

  /**
   * 序列化方式
   */
  contentType: ContentType;

  /**
   * 压缩类型
   */
  contentEncoding: CompressType;

  /**
   * 被调方接入点
   */
  remote: Endpoint;

  /**
   * 透传的上下文
   */
  context: Record<string, Uint8Array>;

  /**
   * 中间件链路耗时上报开关
   */
  needTraceCost: boolean;
  
  clientMiddlewares: Middleware[],
  dispatherMiddlewares: ds.Middleware<any>[],
}

/**
 * option 的默认值
 * 使用 const enum 实现常量，编译时内联
 */
// eslint-disable-next-line @typescript-eslint/naming-convention
export const enum DEFAULT_OPTION {
  TIMEOUT = 5_000,
  CALLER = '.node.client.Default',
  // 窗口大小 默认2^21 - 1
  INIT_WINDOW_SIZE = 2097151,
}

export abstract class BaseContext<CustomOptions extends Options> {
  public remote?: Endpoint;
  public local?: Endpoint;

  public startTime = 0;
  public endTime = 0;
  public middlewareTiming?: [number, number][];

  public error?: RpcError;
  public timeout = false;

  public done = false;
  protected pending = false;

  // eslint-disable-next-line @typescript-eslint/naming-convention
  protected _costTime = 0;
  // eslint-disable-next-line @typescript-eslint/naming-convention
  protected _traceCostTime?: number[];

  protected timer?: ReturnType<typeof setTimeout>;
  protected destroyed = false;

  private deferred?: Deferred<void, RpcError>;
  public abstract type: 'unary' | 'stream';
  public abstract request: /*StreamInitMessage |*/ UnaryRequestMessage;
  public abstract response?: /*StreamInitMessage |*/ UnaryResponsePacket;
  public abstract options: Pick<Options, 'timeout'> & Partial<CustomOptions>;

  /** 便于观察 */
  public inspect() {
    return this.toJSON();
  }

  /**
   * 开始计时
   * 返回一个 deferred promise
   */
  public start(timer: ReturnType<typeof setTimeout>): Promise<void> {
    DEBUGGER('[start]', `pending:${this.pending}`);
    if (this.deferred !== undefined) {
      return this.deferred.promise;
    }

    this.pending = true;
    this.done = false;
    this.deferred = new Deferred<void, RpcError>();
    this.timer = timer;

    return this.deferred.promise;
  }

  /**
   * 结束请求，start后，均通过该方法触发resolve或reject
   * @param err 发生错误
   */
  public end(err?: RpcError) {
    DEBUGGER('[end]', `done:${this.done}`, `pending:${this.pending}`, err);
    if (this.done) return;
    this.done = true;
    this.pending = false;
    this.error = err;

    DEBUGGER('[end]', 'timer', this.timer);
    /* istanbul ignore else */
    if (this.timer !== undefined) {
      clearTimeout(this.timer);
      this.timer = undefined;
    }

    /* istanbul ignore else */
    if (this.deferred !== undefined) {
      this.deferred.resolve();
    }
  }

  public get costTime() {
    let t = this._costTime;
    if (this._costTime === 0 && this.endTime > 0 && this.startTime > 0) {
      t = (this.endTime - this.startTime) * 1000;
      this._costTime = t;
    }
    return t;
  }

  public get traceCostTime(): number[]|undefined {
    // 使用缓存
    if (Array.isArray(this._traceCostTime)) return this._traceCostTime;

    const timing = this.middlewareTiming;
    // 如果记录的时间点少于2，说明中间件执行异常或者没有中间件执行，不再计算耗时
    if (!Array.isArray(timing) || timing.length < 2) return;

    // middlewareTiming记录的是时间点，计算出各点相减的耗时
    let next;
    this._traceCostTime = timing.reduce((acc: number[], current, index) => {
      if (index < timing.length - 1) {
        next = timing[index + 1];
        // 两点相减计算耗时，单位为ms
        return acc.concat((next[0] * 1000 + next[1] / 1e6) - (current[0] * 1000 + current[1] / 1e6));
      }
      return acc;
    }, []);
    return this._traceCostTime;
  }

  protected toJSON() {
    return {
      done: this.done,
      pending: this.pending,
      local: this.local,
      remote: this.remote,
      startTime: this.startTime,
      endTime: this.endTime,
    };
  }
}
