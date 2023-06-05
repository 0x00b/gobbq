import * as net from 'net'
import * as ex from "../../example/exampb/exam"
import * as kcp from 'kcpjs'
import { createTruncator } from '../src/transport/truncator'


// let block = undefined
// if (algorithm && key && iv) {
//     block = new kcp.AesBlock(algorithm, key, iv)
// }

function connect() {
    const socket = kcp.DialWithOptions({
      conv: 255,
      port: 8899,
      host: "localhost",
    //   block,
      dataShards: 10,
      parityShards: 3,
    })
    if (!socket) {
      console.log("dial failed")
      return
    } 
    let i=0
    // const handleData = createTruncator(this.onFrame.bind(this));
    socket
      .on('recv', (data)=>{
        console.log("recv:",i++, data)
      })
      
    for (let index = 0; index < 1000; index++) { 
      socket.write(Buffer.from("xxxxxxxxxx----"+index)) 
    }

}

connect()
