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

func (st *ServerPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	// parse pkt , get header

	hdr := &bbqpb.RequestHeader{}

	codec.GetCodec(bbqpb.ContentType_proto).Unmarshal(pkt.Data(), hdr)

	_ = hdr.Method

	fmt.Println("recv", string(pkt.PacketBody()))

	npkt := codec.NewPacket()
	npkt.WriteBytes([]byte("test"))

	err := pkt.Src.WritePacket(npkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}
