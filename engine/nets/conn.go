package nets

import (
	"errors"
	"io"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

type conn struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	idleTimeout      time.Duration
	lastVisited      time.Time
	PacketHandler    PacketHandler
	opts             *Options
}

func (st *conn) Name() string {
	return "server"
}

func (st *conn) WritePacket(pkt *codec.Packet) error {
	xlog.Println("write header:", pkt.Header.String())
	return st.packetReadWriter.WritePacket(pkt)
}

func (st *conn) Serve() {
	// defer st.Close()
	for {
		// 检查上游是否关闭
		// select {
		// case <-st.ctx.Done():
		// 	return
		// default:
		// }

		if st.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(st.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				st.lastVisited = now
				err := st.rwc.SetReadDeadline(now.Add(st.idleTimeout))
				if err != nil {
					xlog.Println("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		pkt, release, err := st.packetReadWriter.ReadPacket()
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				xlog.Println("transport: tcpconn  EOF ", err)
				// report.TCPTransportReadEOF.Incr() // 客户端主动断开连接
				return
			}
			if e, ok := err.(net.Error); ok && e.Timeout() { // 客户端超过空闲时间没有发包，服务端主动超时关闭
				xlog.Println("transport: tcpconn  Time out ", err)
				// report.TCPTransportIdleTimeout.Incr()
				return
			}
			// report.TCPTransportReadFail.Incr()
			xlog.Println("transport: tcpconn serve ReadFrame fail ", err)
			return
		}
		// report.TCPTransportReceiveSize.Set(float64(len(req)))

		err = st.handle(pkt, release)
		if err != nil {
			xlog.Println("handle failed", err)
		}
	}
}

func (st *conn) handle(pkt *codec.Packet, release codec.ReleasePkt) error {

	defer release()
	// todo chan

	return st.PacketHandler.HandlePacket(pkt)
}
