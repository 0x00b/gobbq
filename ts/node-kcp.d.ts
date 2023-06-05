
declare module 'node-kcp' {
  import { EventEmitter } from 'events'; 
  class KCP extends EventEmitter {
    constructor(conv: number, options?: {
      address?:string,
      port?:number,
    });
    send(buffer: Buffer | string): number;
  }
  export default KCP;
}