/**
 * 接入点
 */
export interface Endpoint {
  /**
   * 主机地址，一般是 IP，也可能是域名
   */
  host: string,
  /**
   * 端口号
   */
  port: number,
  /**
   * 传输协议
   */
  protocol: 'kcp' | 'tcp' | 'ws' | 'wss',
}

export function getEndpointName(e: Endpoint): string {
  return `${e.protocol}://${e.host}:${e.port}`;
}
