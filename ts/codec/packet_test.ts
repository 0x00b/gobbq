import * as packet from "./packet"
import * as bbq from "../../proto/bbq/bbq"

function test() {

    let hdr = bbq.Header.create()
    hdr.Method = "xxxx"
    hdr.RequestId = "1111"

    let hdrBuf = bbq.Header.encode(hdr)

    let hdrBytes = hdrBuf.finish()

    let data = "qqqqqqqq"

    let total = 4 + hdrBytes.length + data.length
    console.log(total, hdrBytes.length, data.length)

    let buf = Buffer.alloc(4)

    buf.writeIntLE(total, 0, 4)
    let pkt = packet.ReadPacket(buf)
    console.log("pkt:", pkt, buf)

    buf = Buffer.alloc(4)
    buf.writeIntLE(4 + hdrBytes.length, 0, 4)
    let pkt2 = packet.ReadPacket(buf)
    console.log("pkt2:", pkt2, buf)

    let pkt3 = packet.ReadPacket(Buffer.from(hdrBytes))
    console.log("pkt3:", pkt3, hdrBytes)

    let pkt4 = packet.ReadPacket(Buffer.from(data))
    console.log("pkt4:", pkt4, data)


    console.log("================")


    let tPkt = new packet.Packet()
    tPkt.Header.RequestId = "123431"
    tPkt.Header.Method = "test"

    tPkt.WriteBody(Buffer.from("112233"))


    let pkt5 = packet.ReadPacket(tPkt.Buffer)
    console.log("pkt4:", pkt4)
    console.log("pkt5:", pkt5)

}

test()