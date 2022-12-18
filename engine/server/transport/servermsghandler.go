package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/bbqpb"
	"github.com/0x00b/gobbq/engine/codec"
)

type ServerPacketHandler struct {
}

func NewServerPacketHandler(ctx context.Context, conn net.Conn) *ServerPacketHandler {
	st := &ServerPacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *ServerPacketHandler) HandlePacket(c context.Context, msg *codec.Packet) error {

	// parse msg , get header

	hdr := &bbqpb.RequestHeader{}

	codec.GetCodec(bbqpb.ContentType_proto).Unmarshal(msg.Data(), hdr)

	_ = hdr.Method

	fmt.Println("recv", string(msg.PacketBody()))

	newmsg := codec.NewPacket()
	newmsg.WriteBytes([]byte("test"))

	err := msg.Src.WritePacket(newmsg)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}
