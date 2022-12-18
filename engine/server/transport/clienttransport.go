package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type ClientTransport struct {
	*conn
}

func NewClientTransport(ctx context.Context, rawConn net.Conn) *ClientTransport {
	ct := &ClientTransport{
		conn: &conn{
			rwc:              rawConn,
			ctx:              context.Background(),
			packetReadWriter: codec.NewPacketReadWriter(context.Background(), rawConn),
			PacketHandler:    NewServerPacketHandler(context.Background(), rawConn),
		},
	}
	// ct.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, ct)
	return ct
}

func (ct *ClientTransport) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	newpkt := codec.NewPacket()
	newpkt.WriteBytes([]byte("test"))

	err := ct.WritePacket(newpkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}

func (ct *ClientTransport) Invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error {

	// req

	// register rsp use req.GetPacketID()

	return nil
}
