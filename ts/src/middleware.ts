import { InitializeOptions } from './client';
import type { BaseContext, Options } from './context'; 

/**
 * 中间件定义
 *
 * @typeParam CustomOptions - 自定义的选项，基于内置选项扩展
 */
export interface Middleware<
  CustomOptions extends Options = Options
> {
  /**
   * 中间件函数的实现
   * @param rpc rpc 上下文
   * @param next 执行下一层中间件
   * - 在 `execute()` 中调用 `await next()` 异步执行下一层中间件，
   *   之后回到 `execute()` 上下中继续执行。
   * - 如果没有调用在 `execute()` 中调用 `next()`，
   *   则在 `execute()` 执行结束后自动执行下一层中间件。
   */
  execute(rpc: BaseContext<CustomOptions>, next: (() => Promise<void>) | (() => void)): Promise<void>;

  /**
   * 中间件可以选择提供一个初始化的方法
   * @param options 包含在初始化阶段确定的选项
   */
  initialize?(options: InitializeOptions): void;

  /**
   * 中间件可以选择提供一个用于销毁的方法
   */
  destroy?(): void;
}

/**
 * 将中间件栈合并为一个中间件
 * @param middlewares 中间件栈
 */
export function compose<
  CustomOptions extends Options
>(mws?: Middleware<CustomOptions>[]): Middleware<CustomOptions> {
  let middlewares: Middleware<CustomOptions>[] = []
  if (mws) {
    middlewares = mws
  }
  if (!Array.isArray(middlewares)) throw new TypeError('Middleware stack must be an array!');

  // eslint-disable-next-line no-restricted-syntax
  for (const mw of middlewares) {
    if (mw === undefined || typeof mw.execute !== 'function') throw new TypeError('Middleware must have execute() method!');
  }

  return {
    initialize(options: InitializeOptions) {
      // eslint-disable-next-line no-restricted-syntax
      for (const mw of middlewares) {
        mw.initialize?.(options);
      }
    },
    execute(rpc: BaseContext<CustomOptions>, next: () => Promise<void>) {
      const { length } = middlewares;
      const { options } = rpc;
      let lastIndex = -1;
      const middlewareTiming: [number, number][] = [];
      // eslint-disable-next-line no-param-reassign
      rpc.middlewareTiming = middlewareTiming;

      const dispatch = options.needTraceCost ? async (currentIndex: number): Promise<void> => {
        if (currentIndex <= lastIndex) throw new Error('next() called multiple times');
        lastIndex = currentIndex;
        if (currentIndex === length) {
          middlewareTiming.push(process.hrtime());
          await next();
          middlewareTiming.push(process.hrtime());
          return;
        }

        // 插件开始执行时间
        middlewareTiming.push(process.hrtime());

        // 执行中间件
        const nextDispatch = dispatch.bind(null, currentIndex + 1);
        await middlewares[currentIndex].execute(rpc, nextDispatch);
        // 如果中间件中没有显式调用 next()，自动调用一次
        if (currentIndex === lastIndex) await nextDispatch();

        // 插件结束执行时间
        middlewareTiming.push(process.hrtime());
      } : async (currentIndex: number): Promise<void> => {
        if (currentIndex <= lastIndex) throw new Error('next() called multiple times');
        lastIndex = currentIndex;
        if (currentIndex === length) return next();

        // 执行中间件
        const nextDispatch = dispatch.bind(null, currentIndex + 1);
        await middlewares[currentIndex].execute(rpc, nextDispatch);
        // 如果中间件中没有显式调用 next()，自动调用一次
        if (currentIndex === lastIndex) return nextDispatch();
      };

      return dispatch(0);
    },
    destroy() {
      // 调用 distroy 的顺序与 initialize 相反
      // eslint-disable-next-line no-restricted-syntax
      for (const mw of middlewares.reverse()) {
        mw.destroy?.();
      }
    },
  };
}
