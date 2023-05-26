export interface Options {
  /**
   * 是否自动重试
   */
  retry: boolean;

  /**
   * 请求超时时间
   */
  timeout: number;

  /**
   * 被调方接入点
   */
  remote: string;
}

export interface Client {
  timeout?: Options['timeout'];
  retry?: Options['retry'];
  services?: Service[]
}

export interface Service {
  timeout?: Options['timeout'];
  retry?: Options['retry'];
  [k: string]: any;
}
