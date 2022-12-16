package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type PacketHandler interface {
	HandlePacket(c context.Context, pkt *codec.Packet) error
}

type ServerTransport struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	ctx              context.Context
	idleTimeout      time.Duration
	lastVisited      time.Time
	PacketHandler    PacketHandler
}

func NewServerTransportWithPacketHandler(ctx context.Context, conn net.Conn, PacketHandler PacketHandler) *ServerTransport {
	st := &ServerTransport{
		rwc:              conn,
		packetReadWriter: codec.NewPacketReadWriter(ctx, conn),
		PacketHandler:    PacketHandler,
		ctx:              ctx,
	}
	return st
}

func NewServerTransport(ctx context.Context, conn net.Conn) *ServerTransport {
	st := &ServerTransport{
		rwc:              conn,
		packetReadWriter: codec.NewPacketReadWriter(ctx, conn),
		ctx:              ctx,
	}
	st.PacketHandler = st
	return st
}

func (st *ServerTransport) Name() string {
	return "server"
}

func (st *ServerTransport) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	newpkt := codec.NewPacket()
	newpkt.WriteBytes([]byte("test"))

	err := st.WritePacket(newpkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}

func (st *ServerTransport) WritePacket(pkt *codec.Packet) error {
	return st.packetReadWriter.WritePacket(pkt)
}

func (st *ServerTransport) parsePacket(ctx context.Context, pkt *codec.Packet) (ci *CallInfo, reqBodyBuff []byte, err error) {

	return
}

func (st *ServerTransport) Serve() {
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

		st.handle(context.Background(), pkt)
	}
}

func (st *ServerTransport) handle(c context.Context, pkt *codec.Packet) error {
	defer pkt.Release()

	switch pkt.GetPacketType() {
	case codec.PacketRPC:
		st.PacketHandler.HandlePacket(c, pkt)
	case codec.PacketSys:
	default:
	}
	return nil
}
