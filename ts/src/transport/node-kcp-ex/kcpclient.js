var kcp = require('node-kcp/build/Release/kcp');
var dgram = require('dgram');

interval = 10

class KCP {
    constructor(conv, { address, port, onrecv, onerror }) {
        this.kcpobj = new kcp.KCP(conv, { address: address, port: port });
        this.kcpobj.stream(1);
        this.kcpobj.nodelay(1, interval, 10, 3);

        var client = dgram.createSocket('udp4');
        client.connect(port, address)

        this.kcpobj.output((data, size, context) => {
            client.send(data, 0, size,/* context.port, context.address*/);
        });

        client.on('error', (err) => {
            client.close();
            if (onerror) {
                onerror(err);
            }
        });

        client.on('message', (data, rinfo) => {
            this.kcpobj.input(data);
            var recv = this.kcpobj.recv();
            if (recv) {
                if (onrecv) {
                    onrecv(recv);
                }
            }
        });

        // todo 优化成check update
        setInterval(() => {
            this.update()
        }, interval);

    }
    update() {
        this.kcpobj.update(Date.now());
    }
    check() {
        return this.kcpobj.check();
    }
    send(data) {
        return this.kcpobj.send(data);
    }
}

exports.KCP = KCP;