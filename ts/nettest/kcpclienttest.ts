import * as kcp from 'kcpjs'
import { buffer } from 'node:stream/consumers'
import * as ex from "../../example/exampb/exam"

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
    console.log('recv:', buff.toString())
})

setInterval(() => {
    let req = ex.SayHelloRequest.create()
    req.text = "test"
    let msg = ex.SayHelloRequest.encode(req)
    let buf = Buffer.from(msg.finish())
    console.log(`send: ${msg}`)
    session.write(buf)
}, 1000)