/**
 *  Node.js Client 实现，提供基础的  客户端功能
 *
 * @remark
 * 错误码定义 {@link ERROR}
 *
 * @remark
 * 传输层 {@link transport}
 *
 * @packageDocumentation
 */

import { Client } from './client';

/**
 * 内置全局 Client 实例
 */
// eslint-disable-next-line @typescript-eslint/naming-convention 
export { Client };
export { ERROR } from './error';

// namespace
export * as transport from './transport';

// type only 
export type { Middleware } from './middleware';
export type { Endpoint } from './endpoint';
export type { BaseContext, /*StreamContext, STREAM_CALL_TYPE, StreamOptions,*/ UnaryContext, Options, UnaryOptions } from './context';

