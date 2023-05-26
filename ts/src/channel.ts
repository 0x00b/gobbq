import { compose } from './middleware';
import { noop } from './utils';
import { getEndpointName } from './endpoint';
import { ERROR, RpcError } from './error';
import { /*StreamContext, STREAM_CALL_TYPE,*/ UnaryContext, DEFAULT_OPTION } from './context';
import { createTransport as builtInCreateTransport } from './transport';
// type only
import type { /*StreamInitMessage,*/ UnaryRequestMessage } from './codec/msg';
import type { Middleware } from './middleware';
import type { Options, UnaryOptions, /*StreamOptions, ClientStreamOptions, ServerStreamOptions, BidiStreamOptions*/ } from './context';
import type { ClientTransport } from './transport';
import type { PluginList } from './plugins';
import type { Endpoint } from './endpoint';
import { Header } from '../../proto/bbq/bbq';


/**
 * 使用 const enum 实现常量，编译时内联
 */
const enum CONST {
  MAX_UINT32 = 2 ** 32 - 1,
}

let id = process.pid;
/**
 * 全局自增的 id
 */
export const genId = () => {
  if (id >= CONST.MAX_UINT32) id = 0;
  return id += 1;
};

/**
 * 初始化时可以指定的选项
 *
 * - `timeout` `callee` `caller` `retry` 为必选
 * - `remote` 为可选
 */
export type InitializeOptions = Pick<Options, 'timeout'> & Partial<Pick<Options, 'remote'>>;

/**
 * {@link Channel} 代表一个逻辑上的调用关系；
 * - {@link Channel} 与  协议中的被调方 `callee` 对应，因此 `callee` 是必需的
 * - 在指定了主调方 `caller` 时，{@link Channel} 与（callee:caller）对应
 * - 一个调用关系通常会部署多个接入点 {@link Endpoint} (protocol://ip:port),
 *   每个接入点分配一个 {@link ClientTransport} 负责传输
 */
export class Channel<CustomOptions extends Options> {
  private readonly options: InitializeOptions;
  /** endpoint name -> transport */
  private readonly transports = new Map<string, ClientTransport>();
  private readonly middleware: Middleware<Options & Partial<CustomOptions>>;
  private destroyed = false;

  private readonly plugins: Pick<PluginList, 'createTransport'>;

  public constructor(
    options: Partial<Pick<Options, 'timeout' | 'remote'>> = {},
    middlewares: Middleware<Options & Partial<CustomOptions>>[] = [],
    plugins: Partial<PluginList> = {},
  ) {
    /* will make sure `callee` `caller` `timeout` `retry` existed */
    this.options = {
      timeout: options.timeout ?? DEFAULT_OPTION.TIMEOUT,
      // caller: options.caller ?? DEFAULT_OPTION.CALLER,
      remote: options.remote,
    };

    this.plugins = {
      createTransport: plugins.createTransport ?? builtInCreateTransport,
    };

    this.middleware = compose<Options & Partial<CustomOptions>>(middlewares);
    this.middleware.initialize!(this.options);
  }

  public toJSON() {
    const { options, plugins, transports } = this;
    return {
      options,
      plugins,
      transports,
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

    this.transports.forEach((transport) => {
      transport.destroy();
    });
    this.transports.clear();
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
    data: Uint8Array,
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
      Body: Buffer.from(data),
    };

    const rpc = new UnaryContext<CustomUnaryOptions>(message, mergedOptions);

    if (opt.remote !== undefined) {
      rpc.remote = opt.remote;
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
          const { response, local } = await transport.addUnary(message);
          console.log("rsp:", response)
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
    let transport = this.transports.get(name);

    if (transport === undefined) {
      console.log('[getTransport]', 'new', name);
      transport = this.plugins.createTransport(remote, noop);
      this.transports.set(name, transport);
    }

    console.log('[getTransport]', 'got', name);
    return transport;
  }
  
}
