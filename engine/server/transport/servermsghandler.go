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

	codec.GetCodec(bbqpb.ContentType_proto).Unmarshal(pkt.PacketBody()[:pkt.GetMsgHeaderLen()], hdr)

	fmt.Println("recv RequestHeader:", hdr.String())
	fmt.Println("recv len:", pkt.GetMsgHeaderLen(), pkt.GetPacketBodyLen())
	fmt.Println("recv data:", string(pkt.PacketBody()[pkt.GetMsgHeaderLen():pkt.GetPacketBodyLen()]))

	_ = hdr.Method

	npkt := codec.NewPacket()
	npkt.WriteBytes([]byte("test"))

	err := pkt.Src.WritePacket(npkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}
