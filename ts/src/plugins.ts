import type { CreateTransport } from './transport';

/**
 * 部分情况下，中间件无法实现的扩展点，通过 Plugins 的形式支持
 */
export interface PluginList {
  createTransport: CreateTransport
};
