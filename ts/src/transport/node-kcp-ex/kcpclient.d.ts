import EventEmitter from "events";


export declare class KCP {
    conv: number;
    port: number;
    host: string;

    constructor(conv: number, options: { address: string, port: number, onrecv?: (data: Buffer) => void, onerror?: (err: any) => void });
    send(b: Buffer | string): number;
    close(): void;
    // setNoDelay(nodelay: number, interval: number, resend: number, nc: number): void;
    check(): void;
    update(): void;
}