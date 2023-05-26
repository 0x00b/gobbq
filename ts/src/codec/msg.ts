import * as bbq from "../../../proto/bbq/bbq"
import { Packet } from "./packet"

//   StreamMessage,
//   StreamInitMessage,

const MaxBufferCap = 16 * 1024 * 1024 //16M

export class UnaryRequestMessage {
  Header: bbq.Header
  Body: Buffer

  constructor(Header: bbq.Header, Body: Buffer) {
    this.Header = Header
    this.Body = Body
  }
}

export interface UnaryResponseMessage extends Packet {
}


// appends slice of bytes to the end of packetBody
export function encode(request: UnaryRequestMessage): Buffer {
  if (!request.Header) {
    throw "invalid req header"
  }

  const header = bbq.Header.encode(request.Header)
  let hdrBytes = header.finish()
  let HeaderLen = 4 + hdrBytes.length
  let TotalLen = HeaderLen + request.Body.length

  // TotalLen 本身要在最前面占4个字节
  const data = Buffer.alloc(4 + TotalLen)

  data.writeInt32LE(TotalLen)
  data.writeInt32LE(HeaderLen, 4)
  data.set(hdrBytes, 8)
  data.set(request.Body, 4 + HeaderLen)

  return data
}

export function decode(tempBuff: Buffer): Packet | null {

  // console.log(`packet.Buffer.length: ${tempBuff.length}`)

  if (tempBuff.length < 4) {
    return null
  }

  let packetDataSize: number
  packetDataSize = tempBuff.readInt32LE()

  console.log(`total len: ${packetDataSize}`)

  if (packetDataSize > MaxBufferCap) {
    console.log(`err1 ${packetDataSize} > ${MaxBufferCap}`)
    return null
  }


  if (packetDataSize + 4 > tempBuff.length) {
    console.log(`err2 ${packetDataSize} > ${packetDataSize + 4}`)
    return null
  }

  let packet = new Packet()

  packet.Buffer = tempBuff

  packet.TotalLen = packetDataSize

  packet.HeaderLen = packet.Buffer.readInt32LE(4)
  // console.log(`packet.HeaderLen: ${packet.HeaderLen}`)
  let packetHeaderData: Uint8Array = packet.Buffer.subarray(4 + 4, 4 + packet.HeaderLen)

  // console.log(`packetHeaderData: ${packetHeaderData}`)

  packet.Header = bbq.Header.decode(packetHeaderData)

  // console.log(`recv raw: ${packet.Buffer.toString()}`)

  return packet
}
