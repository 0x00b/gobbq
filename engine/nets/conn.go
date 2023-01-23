package nets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type conn struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	ctx              context.Context
	idleTimeout      time.Duration
	lastVisited      time.Time
	PacketHandler    PacketHandler
	opts             *Options
}

func (st *conn) Name() string {
	return "server"
}

func (st *conn) WritePacket(pkt *codec.Packet) error {
	return st.packetReadWriter.WritePacket(pkt)
}

func (st *conn) Serve() {
	// defer st.Close()
	for {
		// 检查上游是否关闭
		select {
		case <-st.ctx.Done():
			return
		default:
		}

		if st.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(st.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				st.lastVisited = now
				err := st.rwc.SetReadDeadline(now.Add(st.idleTimeout))
				if err != nil {
					fmt.Println("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		pkt, err := st.packetReadWriter.ReadPacket()
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				// report.TCPTransportReadEOF.Incr() // 客户端主动断开连接
				return
			}
			if e, ok := err.(net.Error); ok && e.Timeout() { // 客户端超过空闲时间没有发包，服务端主动超时关闭
				// report.TCPTransportIdleTimeout.Incr()
				return
			}
			// report.TCPTransportReadFail.Incr()
			fmt.Println("transport: tcpconn serve ReadFrame fail ", err)
			return
		}
		// report.TCPTransportReceiveSize.Set(float64(len(req)))

		err = st.handle(context.Background(), pkt)
		if err != nil {
			fmt.Println("handle failed", err)
		}
	}
}

func (st *conn) handle(c context.Context, pkt *codec.Packet) error {
	defer pkt.Release()

	return st.PacketHandler.HandlePacket(c, pkt)
}
