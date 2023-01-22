package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
)

var _ server.PacketHandler = &GatePacketHandler{}

type GatePacketHandler struct {
}

func NewGatePacketHandler() *GatePacketHandler {
	st := &GatePacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *GatePacketHandler) HandlePacket(c context.Context, opts *server.ServerOptions, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	// hdr := &proto.Header{}

	hdr := pkt.GetHeader()

	// codec.DefaultCodec.Unmarshal(pkt.PacketBody()[:pkt.GetMsgHeaderLen()], hdr)

	// fmt.Println("recv RequestHeader:", hdr.String())
	// fmt.Println("recv len:", pkt.GetMsgHeaderLen(), pkt.GetPacketBodyLen())
	fmt.Println("recv data:", string(pkt.PacketBody()))

	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		fmt.Println("err mothod")
		return errors.New("err mothod")
	}
	service := sm[:pos]
	method := sm[pos+1:]

	_ = service
	_ = method

	// send to dispather
	return nil
}
