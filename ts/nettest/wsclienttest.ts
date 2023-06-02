import * as net from 'net'
import * as ex from "../../example/exampb/exam"

let block = undefined
// if (algorithm && key && iv) {
//     block = new kcp.AesBlock(algorithm, key, iv)
// }

function connect() {
    var serverAddr = 'ws://1.12.243.253/ws'
    console.log("正在连接 " + serverAddr + ' ...')
    var websocket = new WebSocket(serverAddr)
    this.websocket = websocket

    websocket.binaryType = 'arraybuffer'
    console.log(websocket)
    var gameclient = this

    //连接发生错误的回调方法
    websocket.onerror = function () {
        console.log("WebSocket连接发生错误")
    }

    //连接成功建立的回调方法
    websocket.onopen = function () {
        console.log("WebSocket连接成功")
        let req = ex.SayHelloRequest.create()
        req.text = "test"
        let msg = ex.SayHelloRequest.encode(req)
        let buf = Buffer.from(msg.finish())
        console.log(`send: ${msg}`)
        websocket.send(buf)
    }

    //接收到消息的回调方法
    websocket.onmessage = function (event) {
        var data = event.data
        console.log("收到数据：", typeof (data), data.length)
        gameclient.onRecvData(data)
    }

    //连接关闭的回调方法
    websocket.onclose = function () {
        console.log("WebSocket连接关闭")
    }

    //监听窗口关闭事件，当窗口关闭时，主动去关闭websocket连接，防止连接还没断开就关闭窗口，server端会抛异常。
    window.onbeforeunload = function () {
        console.log("onbeforeunload")
    }
}

connect()
