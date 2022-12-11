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

type StreamServer struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	ctx              context.Context
	idleTimeout      time.Duration
	lastVisited      time.Time
}

func NewStreamServer(ctx context.Context, conn net.Conn) *StreamServer {
	return &StreamServer{
		rwc:              conn,
		packetReadWriter: codec.NewPacketReadWriter(ctx, conn),
		ctx:              ctx,
	}
}

func (c *StreamServer) Close() {
}

func (c *StreamServer) Serve() {
	defer c.Close()
	for {
		// 检查上游是否关闭
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if c.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(c.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				c.lastVisited = now
				err := c.rwc.SetReadDeadline(now.Add(c.idleTimeout))
				if err != nil {
					fmt.Println("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		pkt, err := c.packetReadWriter.ReadPacket()
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

		c.handle(pkt)
	}
}

func (c *StreamServer) handle(packet *codec.Packet) {
	defer packet.Release()

	switch packet.GetPacketType() {
	case codec.PacketRPC:
		c.handleRPC(packet)
	case codec.PacketPing:
	default:
	}

}

func (c *StreamServer) handleRPC(packet *codec.Packet) {

	fmt.Println("recv", string(packet.PacketBody()))

	newpkt := codec.NewPacket()
	newpkt.WriteBytes([]byte("test"))

	err := c.packetReadWriter.WritePacket(newpkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	// entityID := packet.EntityID()
	// _ = entityID

	// sm := packet.Method()
	// if sm != "" && sm[0] == '/' {
	// 	sm = sm[1:]
	// }
	// pos := strings.LastIndex(sm, "/")
	// if pos == -1 {
	// 	return
	// }
	// service := sm[:pos]
	// method := sm[pos+1:]

	// srv, knownService := s.services[service]
	// if knownService {
	// 	if md, ok := srv.methods[method]; ok {
	// 		s.processUnaryRPC(t, packet, srv, md)
	// 		return
	// 	}
	// }
}
