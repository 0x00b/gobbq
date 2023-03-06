import * as kcp from 'kcpjs'
import { buffer } from 'node:stream/consumers'
import * as ex from "../../example/exampb/exam"
import { Packet, ReadPacket } from '../codec/packet'

let block = undefined
// if (algorithm && key && iv) {
//     block = new kcp.AesBlock(algorithm, key, iv)
// }

// client
const session = kcp.DialWithOptions({
    conv: 255,
    port: 8899,
    host: "127.0.0.1",
    block,
    dataShards: 10,
    parityShards: 3,
})

session.on("recv", (buff: Buffer) => {

    let pkt = ReadPacket(buff)
    if (pkt == null) {
        console.log('recv null xx', buff)
        return
    }

    console.log('recv:', JSON.stringify(pkt))
})

setInterval(() => {

    let pkt = new Packet()

    pkt.Header.Method = "xxxx"
    pkt.Header.RequestId = "1111"

    let req = ex.SayHelloRequest.create()
    req.text = "test"

    let msg = ex.SayHelloRequest.encode(req)

    let buf = Buffer.from(msg.finish())

    pkt.WriteBody(buf)

    console.log(`send: ${JSON.stringify(pkt)}`)

    session.write(pkt.Buffer)

}, 1000)