// import { configure } from '@/-config-env';
import { DEFAULT_OPTION } from './context';
import { Channel } from './channel';
// type only
import type { Client as ClientConfigure, Service as ServiceConfigure } from './configure';
import type { Options } from './context';
import type { Middleware } from './middleware';
import type { PluginList } from './plugins';
import { ServiceDefinition } from './dispatcher/service';
import { Dispatcher } from './dispatcher/dispatcher';
import { ClientTransport } from './transport';

/**
 * client 仅用于存储一些 “全局” 的配置项（timeout/caller）
 * 并提供工厂方法来创建 Channel
 */
export class Client {

  public channel: Channel<any>;
  private dispather: Dispatcher<any>;
  // private transport: ClientTransport;
  
  public constructor(options: Partial<Pick<Options, 'timeout'>> = Object.create(null)) {
    // this.options = {
    // timeout: options.timeout ?? this.configure.timeout ?? DEFAULT_OPTION.TIMEOUT,
    // caller: options.caller ?? this.configure.caller ?? DEFAULT_OPTION.CALLER,
    // };


    // this.transport = undefined

    this.channel = this.createChannel("")

    this.dispather = new Dispatcher("")

  }

  public toJSON() {
    return {
      // configure: this.configure,
      // options: this.options,
    };
  }

  /**
   * 用于创建 Channel 的工厂方法 (Factory Method)
   * @param callee 被调方标识
   * @param middlewares 中间件
   * @param options 支持在自定义参数，会在中间件中透传
   */
  private createChannel(
    callee: string,
    options: Partial<Pick<Options, 'timeout' | 'remote'>> = Object.create(null),
    middlewares: Middleware[] = [],
    plugins: Partial<PluginList> = {},
  ) {
    // const config: ServiceConfigure | undefined = this.configure.services?.find(service => service.callee === callee);

    /** 处理选项的优先级
     * `callee` 作为必选参数，优先级最高
     * `options` 为工厂方法参数，优先级次之
     * `config` 为配置，优先级低于 `options`
     * `this.options` 为全局配置项，优先级最低
     * 同时保证 `callee` `caller` `timeout` 不是 `undefined` */
    const mergedOptions = {
      // ...this.options,
      // ...config,
      ...options,
      callee,
    };

    return new Channel(mergedOptions, middlewares, plugins);
  }

  // public RegisterClientImpl(def:ServiceDefinition, impl:any){

  // }

}
