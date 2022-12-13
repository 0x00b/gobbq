package stream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type ClientTransport struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	ctx              context.Context
	idleTimeout      time.Duration
	lastVisited      time.Time
}

func NewClientTransport(ctx context.Context, conn net.Conn) *ClientTransport {
	ct := &ClientTransport{
		ctx:              ctx,
		rwc:              conn,
		packetReadWriter: codec.NewPacketReadWriter(ctx, conn),
	}
	return ct
}

func (c *ClientTransport) Close() {
}

func (ct *ClientTransport) WritePacket(pkt *codec.Packet) error {
	return ct.packetReadWriter.WritePacket(pkt)
}

func (ct *ClientTransport) Serve() {

	defer ct.Close()
	for {
		// 检查上游是否关闭
		select {
		case <-ct.ctx.Done():
			return
		default:
		}

		if ct.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(ct.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				ct.lastVisited = now
				err := ct.rwc.SetReadDeadline(now.Add(ct.idleTimeout))
				if err != nil {
					fmt.Println("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		pkt, err := ct.packetReadWriter.ReadPacket()
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				// report.TCPServerTransportReadEOF.Incr() // 客户端主动断开连接
				return
			}
			if e, ok := err.(net.Error); ok && e.Timeout() { // 客户端超过空闲时间没有发包，服务端主动超时关闭
				// report.TCPServerTransportIdleTimeout.Incr()
				return
			}
			// report.TCPServerTransportReadFail.Incr()
			fmt.Println("transport: tcpconn serve ReadFrame fail ", err)
			return
		}
		// report.TCPServerTransportReceiveSize.Set(float64(len(req)))

		ct.handle(pkt)
	}
}

func (st *ClientTransport) handle(packet *codec.Packet) {
	defer packet.Release()

	switch packet.GetPacketType() {
	case codec.PacketRPC:
		st.handleRPC(packet)
	case codec.PacketPing:
	default:
	}

}

func (ct *ClientTransport) handleRPC(packet *codec.Packet) {

	fmt.Println("recv", string(packet.PacketBody()))

	// newpkt := codec.NewPacket()
	// newpkt.WriteBytes([]byte("test"))

	// err := ct.packetReadWriter.WritePacket(newpkt)
	// if err != nil {
	// 	fmt.Println("WritePacket", err)
	// }
}
