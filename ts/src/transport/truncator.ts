type PacketHandler = (packet: Buffer) => void;

type Truncator = (chunk: Buffer) => number;

/**
 * 用于从字节流中，切分出完整数据包
 * @param handler { PacketHandler } 处理完整数据包的回调
 */
export function createTruncator(handler: PacketHandler): Truncator {
  let chunks: Buffer[] = [];
  let length = 0;

  let packetLength = 4;

  return (chunk: Buffer): number => {
    chunks.push(chunk);

    console.log("chunk.length:", chunk.length)

    length += chunk.length;

    // 不足一个 packet，等待下次调用
    if (length < packetLength) {
      return length;
    }

    let buf = Buffer.concat(chunks, length);

    do {
      // 未读取过包长度
      if (packetLength === 4) {
        packetLength = buf.readInt32LE();
        if (packetLength < 4) {
          throw new Error(`invalid packet length ${packetLength}`);
        }
      }

      // 不足一个 packet
      if (buf.length + 4 < packetLength) {
        break;
      }

      // console.log("recv packetLength:", packetLength)

      handler(buf.subarray(0, packetLength + 4));

      buf = buf.subarray(packetLength + 4);

      // console.log("reserve buf:", buf)

      packetLength = 4;
    } while (buf.length >= 4);

    // 剩余 buffer
    chunks = [buf];
    length = buf.length;

    // 返回剩余 buffer
    return length;
  };
}

// let handle = createTruncator((buffer:Buffer)=>{
//   console.log("res:",buffer)
// })

// let temp =Buffer.from("\x0a\x00\x00\x00123456\x0a\x00\x00\x007890") 
// handle(temp)