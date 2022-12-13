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

type ServerTransport struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	ctx              context.Context
	idleTimeout      time.Duration
	lastVisited      time.Time
}

func NewServerTransport(ctx context.Context, conn net.Conn) *ServerTransport {
	return &ServerTransport{
		rwc:              conn,
		packetReadWriter: codec.NewPacketReadWriter(ctx, conn),
		ctx:              ctx,
	}
}

func (c *ServerTransport) Close() {
}

func (st *ServerTransport) Serve() {
	defer st.Close()
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

		st.handle(pkt)
	}
}

func (st *ServerTransport) handle(packet *codec.Packet) {
	defer packet.Release()

	switch packet.GetPacketType() {
	case codec.PacketRPC:
		st.handleRPC(packet)
	case codec.PacketPing:
	default:
	}

}

func (st *ServerTransport) handleRPC(packet *codec.Packet) {

	fmt.Println("recv", string(packet.PacketBody()))

	newpkt := codec.NewPacket()
	newpkt.WriteBytes([]byte("test"))

	err := st.packetReadWriter.WritePacket(newpkt)
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
