/**
 * Compose `middleware` returning
 * a fully valid middleware comprised
 * of all those which are passed.
 *
 * @param {Array} middleware
 * @return {Function}
 * @api public
 */

export type Next = () => Promise<any>;
export type Middleware<T> = (context: T, next: Next) => any;
export type ComposedMiddleware<T> = (context: T, next?: Next) => Promise<void>;

/**
 * koa-compose
 * https://github.com/koajs/compose/blob/master/index.js
 * @param middleware
 */
export function compose<T>(middleware: Middleware<T>[]): ComposedMiddleware<T> {
  if (!Array.isArray(middleware)) {
    throw new TypeError('Middleware stack must be an array!');
  }

  // eslint-disable-next-line no-restricted-syntax
  for (const fn of middleware) {
    if (typeof fn !== 'function') {
      throw new TypeError('Middleware must be composed of functions!');
    }
  }

  /**
   * @param {Object} context
   * @return {Promise}
   * @api public
   */
  return function (context: T, next?: Next): Promise<void> {
    // last called middleware #
    let index = -1;
    return dispatch(0);

    function dispatch(i: number): Promise<any> {
      if (i <= index) {
        return Promise.reject(new Error('next() called multiple times'));
      }

      index = i;
      const fn = i === middleware.length ? next : middleware[i];

      if (!fn) {
        return Promise.resolve();
      }

      try {
        return Promise.resolve(fn(context, dispatch.bind(null, i + 1)));
      } catch (err) {
        return Promise.reject(err);
      }
    }
  };
}
