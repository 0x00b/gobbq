import * as net from 'net'
import * as ex from "../../example/exampb/exam"

let block = undefined
// if (algorithm && key && iv) {
//     block = new kcp.AesBlock(algorithm, key, iv)
// }

// client
const client = net.createConnection({
    port: 8899,
    host: "127.0.0.1",
})

client.on('connect', () => {
    // 向服务器发送数据
    client.write('Nodejs 技术栈')

    setTimeout(() => {
        let req = ex.SayHelloRequest.create()
        req.text = "test"
        let msg = ex.SayHelloRequest.encode(req)
        let buf = Buffer.from(msg.finish())
        console.log(`send: ${msg}`)
        client.write(buf)
    }, 1000)
})

client.on('data', buffer => {
    console.log(buffer.toString())
})

// 例如监听一个未开启的端口就会报 ECONNREFUSED 错误
client.on('error', err => {
    console.error('服务器异常：', err)
})

client.on('close', err => {
    console.log('客户端链接断开！', err)
}) 