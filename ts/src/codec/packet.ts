import * as bbq from "../../../proto/bbq/bbq"


export class Packet {
  Header: bbq.Header
  TotalLen: number
  HeaderLen: number
  Buffer: Buffer //totallen(4) + headerlen(4) + header + body

  constructor() {
    this.Header = bbq.Header.create()
    this.TotalLen = 0
    this.HeaderLen = 0
    this.Buffer = Buffer.alloc(0)
  }

  // returns the total packetBody of packet
  public PacketBody(): Buffer {
    return this.Buffer?.subarray(4+this.HeaderLen, 4+this.TotalLen)
  }

}
