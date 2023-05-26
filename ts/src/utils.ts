export function required(message: string): never {
  throw new TypeError(`required: ${message}`);
}

export const noop = () => { /* Do nothing */ };

export class Deferred<T, E> {
  public resolve: (res: T) => void = noop;
  public reject: (err: E) => void = noop;

  // eslint-disable-next-line @typescript-eslint/naming-convention
  private _promise: Promise<T>;

  public constructor() {
    this._promise = new Promise<T>((resolve, reject) => {
      this.resolve = resolve;
      this.reject = reject;
    });
  }

  public get promise() {
    return this._promise;
  }
}
