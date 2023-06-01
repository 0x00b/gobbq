// import { configure } from '@/-config-env';
import { DEFAULT_OPTION, UnaryContext, UnaryOptions } from './context';
// type only
import type { Client as ClientConfigure, Service as ServiceConfigure } from './configure';
import type { Options } from './context';
import { compose, Middleware } from './middleware';
import type { PluginList } from './plugins';
import { ServiceDefinition } from './dispatcher/service';
import { Dispatcher } from './dispatcher/dispatcher';
import { ClientTransport, createTransport, UnaryResult } from './transport';
import { Deferred, noop } from './utils';
import { ERROR, RpcError } from './error';

import { Header, RequestType } from '../../proto/bbq/bbq';
import { encode, UnaryRequestMessage } from './codec/msg';
import { Endpoint, getEndpointName } from './endpoint';
import { Packet } from './codec/packet';

/**
 * 初始化时可以指定的选项
 *
 * - `timeout` `callee` `caller` `retry` 为必选
 * - `remote` 为可选
 */
export type InitializeOptions = Pick<Options, 'timeout'> & Partial<Pick<Options, 'remote'>>;

/**
 * client 仅用于存储一些 “全局” 的配置项（timeout/caller）
 * 并提供工厂方法来创建 Channel
 */
export class Client<CustomOptions extends Options> {

  private readonly options: InitializeOptions;
  /** endpoint name -> transport */
  private transport: ClientTransport;
  private readonly middleware: Middleware<Options & Partial<CustomOptions>>;
  private destroyed = false;

  private readonly plugins: Pick<PluginList, 'createTransport'>;

  private dispather: Dispatcher<any>;
  // private transport: ClientTransport;

  /** 未收到回包的 unary  */
  private deferredUnary = new Map<string, Deferred<UnaryResult, RpcError>>();

  public constructor(
    def: ServiceDefinition, impl: any, 
    options: Partial<Options> = Object.create(null),
  ) {

    this.dispather = new Dispatcher(def, impl)
    options.dispatherMiddlewares?.forEach(m => {
      this.dispather.use(m)
    })

    /* will make sure `callee` `caller` `timeout` `retry` existed */
    this.options = {
      timeout: options.timeout ?? DEFAULT_OPTION.TIMEOUT,
      // caller: options.caller ?? DEFAULT_OPTION.CALLER,
      remote: options.remote,
    };

    this.plugins = {
      createTransport: createTransport,
    };

    if (!this.options.remote) {
      throw new Error("this.options.remote is undefined")
    }

    this.transport = this.plugins.createTransport(this.options.remote, noop, this.onUnaryMessage.bind(this));

    this.middleware = compose<Options & Partial<CustomOptions>>(options.clientMiddlewares);
    this.middleware.initialize!(this.options);

  }

  public toJSON() {
    const { options, plugins, transport } = this;
    return {
      options,
      plugins,
      transport,
      // configure: this.configure,
      // options: this.options,
    };
  }

  /**
   * 销毁
   */
  public destroy() {
    console.log('[destroy]', `destroyed:${this.destroyed}`);
    /* istanbul ignore if */
    if (this.destroyed) return;
    this.destroyed = true;

    this.transport.destroy();
    this.clear();
    this.middleware.destroy!();
  }

  /**
   * 普通Rpc调用
   * @param Header 接口描述
   * @param data 请求二进制数据
   * @param opt
   */
  public async unaryInvoke<CustomUnaryOptions extends CustomOptions & UnaryOptions>(
    Header: Header,
    data: Buffer,
    opt: Partial<CustomUnaryOptions> = {},
  ): Promise<UnaryContext<CustomUnaryOptions>> {
    console.log('[unaryInvoke]', Header.Method, opt);

    const mergedOptions = {
      ...this.options,
      retry: false,
      callType: 0 as const,
      ...opt,
    };

    const message: UnaryRequestMessage = {
      Header: Header,
      Body: data,
    };

    const rpc = new UnaryContext<CustomUnaryOptions>(message, mergedOptions);

    if (this.options.remote !== undefined) {
      rpc.remote = this.options.remote;
    }

    // eslint-disable-next-line @typescript-eslint/naming-convention
    let connecting = false;
    setImmediate(() => {
      this.middleware.execute(rpc, async () => {
        if (rpc.done) return;

        const { remote, request } = rpc;
        console.log('[unaryInvoke]', `remote:${JSON.stringify(remote)}`);

        if (!remote) {
          rpc.end(new RpcError(ERROR.CLIENT_ROUTER_ERR, 'no remote'));
          return;
        }

        const name = getEndpointName(remote);
        // /* istanbul ignore else */
        // if (request.context) {
        //   request.context['-remote'] = Buffer.from(name);
        // }

        const transport = this.getTransport(remote, name);

        try {
          connecting = true;
          await transport.connect();
        } catch (error) {
          rpc.end(new RpcError(ERROR.CLIENT_CONNECT_ERR, "error"));
          return;
        }

        if (rpc.done) return;

        try {
          rpc.startTime = process.uptime();
          // const { response, local } = await transport.addUnary(message);
          console.log("req:", message)
          const { response, local } = await this.addUnary(message);
          console.log("22rsp:", response)
          rpc.endTime = process.uptime();
          rpc.respond(response, local);
        } catch (error) {
          /* istanbul ignore if */
          if (error instanceof RpcError) {
            rpc.end(error);
          } else {
            rpc.end(new RpcError(ERROR.INVOKE_UNKNOWN_ERR, "error"));
          }
        }
      }).then(() => {
        rpc.end(); // 中间件执行结束后，才结束 RPC
      }, (error) => {
        /* istanbul ignore if */
        if (error instanceof RpcError) {
          rpc.end(error);
        } else {
          rpc.end(new RpcError(ERROR.INVOKE_UNKNOWN_ERR, error.message, error.cause ?? error));
        }
      });
    });

    console.log(mergedOptions.timeout)
    const timeout = setTimeout(() => {
      /* istanbul ignore if */
      if (rpc.done) return;
      rpc.timeout = true;

      // 未连接
      if (!connecting) {
        rpc.end(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, 'Timeout before connecting'));
        return;
      }

      // 未发送
      if (rpc.startTime === 0) {
        rpc.end(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, 'Timeout before sending'));
        return;
      }

      // 未返回
      if (rpc.endTime === 0) {
        rpc.endTime = process.uptime();
        rpc.end(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, 'Timeout before received'));
        return;
      }

      rpc.end(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, 'Timeout after received'));
    }, mergedOptions.timeout);

    await rpc.start(timeout);
    return rpc;
  }

  /**
   * 获取一个 { @license Endpoint } 对应的 { @link ClientTransport }
   * @param remote 远程接入点
   * @param name 接入点名
   */
  private getTransport(remote: Endpoint, name: string): ClientTransport {
    // let transport = this.transports.get(name);

    if (this.transport === undefined) {
      console.log('[getTransport]', 'new', name);
      this.transport = this.plugins.createTransport(remote, noop, this.onUnaryMessage);
      // this.transport.set(name, transport);
    }

    console.log('[getTransport]', 'got', name);
    return this.transport;
  }


  /**
   * 主动删除一个 unary 请求，会产生一个 RPC Error
   * @param requestId 请求 id
   * @param err 指定 RPC Error
   */
  public removeUnary(
    requestId: string,
    err = new RpcError(ERROR.CLIENT_CANCELED_ERR, 'unary request removed'),
  ) {
    const deferred = this.deferredUnary.get(requestId);
    if (deferred === undefined) return;
    this.deferredUnary.delete(requestId);
    deferred.reject(err);
  }

  private onUnaryMessage(pkt: Packet) {

    // request
    if (pkt.Header.RequestType == RequestType.RequestRequest) {
      this.dispather.onUnaryMessage(pkt, this.transport)
      return
    }

    // response
    const { RequestId } = pkt.Header;
    const deferred = this.deferredUnary.get(RequestId);
    /* istanbul ignore if */
    if (deferred === undefined) return;
    this.deferredUnary.delete(RequestId);
    deferred.resolve({
      response: pkt,
      local: this.transport.local(),
      remote: this.transport.remote(),
    }
    );
  }

  private clear() {
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
        d.reject(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, "clear"));
      });
      this.deferredUnary.clear();
    }
  }

  /**
   * 添加一个 unary 请求
   * @returns 返回 Promise
   * - 收到来自 server 的响应时，Promise 返回 UnaryRequestMessage
   * - 发送过程中发生异常时，Promise 抛出一个 RpcError
   * @param req 请求 Message
   */
  public async addUnary(req: UnaryRequestMessage): Promise<UnaryResult> {
    console.log('[add] unary', `requestId:${req.Header.RequestId}`);

    let buf: Buffer;
    try {
      buf = encode(req);
    } catch (error) {
      throw new RpcError(ERROR.CLIENT_ENCODE_ERR, "error.message, error");
    }

    const deferred = new Deferred<UnaryResult, RpcError>();
    this.deferredUnary.set(req.Header.RequestId, deferred);

    // console.log(req.Header)

    setTimeout(() => {
      this.deferredUnary.delete(req.Header.RequestId);
      deferred.reject(new RpcError(ERROR.CLIENT_INVOKE_TIMEOUT_ERR, 'Timeout'));
    }, req.Header.Timeout);

    // 稍后发送
    setImmediate(() => {
      try {
        console.log("send:", buf)
        this.transport.send(buf);
      } catch (error) {
        this.deferredUnary.delete(req.Header.RequestId);
        deferred.reject(new RpcError(ERROR.CLIENT_NETWORK_ERR, "error.message, error"));
      }
    });

    return deferred.promise;
  }

}
