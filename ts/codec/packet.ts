import * as bbq from "../../proto/bbq/bbq"


const MaxBufferCap = 16 * 1024 * 1024 //16M

export class Packet {
  Header: bbq.Header
  TotalLen: number
  HeaderLen: number
  Buffer: Buffer

  constructor() {
    this.Header = bbq.Header.create()
    this.TotalLen = 0
    this.HeaderLen = 0
    this.Buffer = Buffer.alloc(0)
  }

  // appends slice of bytes to the end of packetBody
  public WriteBody(b: Buffer): void {
    if (!this.Header) {
      return
    }
    const header = bbq.Header.encode(this.Header)
    let hdyBytes = header.finish()
    this.HeaderLen = 4 + hdyBytes.length
    this.TotalLen = this.HeaderLen + b.length

    // TotalLen 本身要在最前面占4个字节
    const data = Buffer.alloc(4 + this.TotalLen)

    data.writeInt32LE(this.TotalLen)
    data.writeInt32LE(this.HeaderLen, 4)
    data.set(hdyBytes, 8)
    data.set(b, 4 + this.HeaderLen)

    this.Buffer = data
  }

  // returns the total packetBody of packet
  public PacketBody(): Uint8Array | null {
    if (!this.Buffer) {
      return null
    }
    return this.Buffer.subarray(this.HeaderLen, this.TotalLen)
  }

}

let packet = new Packet()

export function ReadPacket(tempBuff: Buffer): Packet | null {

  packet.Buffer = Buffer.concat([packet.Buffer, tempBuff])

  // console.log(`packet.Buffer.length: ${packet.Buffer.length}`)

  if (packet.Buffer.length < 4) {
    return null
  }

  let packetDataSize: number
  // packetDataSize = new DataView(packet.Buffer.buffer).getInt32(0, true)
  packetDataSize = packet.Buffer.readInt32LE()
  // console.log(`total len: ${packetDataSize}`, packet.Buffer) 

  if (packetDataSize > MaxBufferCap) {
    return null
  }

  if (packetDataSize + 4 > packet.Buffer.length) {
    return null
  }

  packet.TotalLen = packetDataSize

  // packet.HeaderLen = new DataView(packet.Buffer.buffer).getInt32(4, true)
  packet.HeaderLen = packet.Buffer.readInt32LE(4)
  // console.log(`packet.HeaderLen: ${packet.HeaderLen}`)
  let packetHeaderData: Uint8Array = packet.Buffer.subarray(4 + 4, 4 + packet.HeaderLen)

  // console.log(`packetHeaderData: ${packetHeaderData}`)

  packet.Header = bbq.Header.decode(packetHeaderData)

  // console.log(`recv raw: ${packet.Buffer.toString()}`)

  let pkt = packet

  packet = new Packet()
  packet.Buffer = packet.Buffer.subarray(4 + packet.TotalLen)

  // pc.readMsgCnt++

  return pkt
}
