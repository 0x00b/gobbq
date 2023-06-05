
import { KCP } from "../src/transport/node-kcp-ex/kcpclient";


 
const kcp = new KCP(1234, {
    address:"localhost",
    port:8899,
    onrecv:(data:Buffer)=>{
      console.log("recv:",data.toString())
    }
});

// kcp.on('data', (data: Buffer) => {
//   console.log(`Received data: ${data.toString()}`);
// });

kcp.send('Hello, world!');


// var kcp = require('node-kcp');
// var kcpobj = new kcp.KCP(123, {address: '127.0.0.1', port: 41234});
// var dgram = require('dgram');
// var client = dgram.createSocket('udp4');
// var msg = 'hello world';
// var idx = 1;
// var interval = 200;

// kcpobj.nodelay(0, interval, 0, 0);

// kcpobj.output((data, size, context) => {
//     client.send(data, 0, size, context.port, context.address);
// });

// client.on('error', (err) => {
//     console.log(`client error:\n${err.stack}`);
//     client.close();
// });

// client.on('message', (msg, rinfo) => {
//     kcpobj.input(msg);
// });

// setInterval(() => {
//     kcpobj.update(Date.now());
//     var recv = kcpobj.recv();
//     if (recv) {
//         console.log(`client recv ${recv}`);
//         kcpobj.send(msg+(idx++));
//     }
// }, interval);

// kcpobj.send(msg+(idx++));
