package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/proto"
)

var _ server.PacketHandler = &GatePacketHandler{}

type GatePacketHandler struct {
}

func NewGatePacketHandler() *GatePacketHandler {
	st := &GatePacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *GatePacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

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
	// ed := st.opts.Entities[service]
	// mt := ed.Methods[method]
	// dec := func(v interface{}) error {
	// 	reqbuf := pkt.PacketBody()
	// 	err := codec.GetCodec(proto.ContentType(hdr.GetContentType())).Unmarshal(reqbuf, v)
	// 	return err
	// }

	// rsp, err := mt.Handler(ed.ServiceImpl, c, dec, nil)

	npkt := codec.NewPacket()

	rhdr := &proto.Header{
		Version:      hdr.Version,
		RequestId:    hdr.RequestId,
		Timeout:      hdr.Timeout,
		Method:       hdr.Method,
		TransInfo:    hdr.TransInfo,
		ContentType:  hdr.ContentType,
		CompressType: hdr.CompressType,
	}
	npkt.SetHeader(rhdr)

	rbyte, err := codec.DefaultCodec.Marshal(rhdr)
	if err != nil {
		fmt.Println("WritePacket", err)
		return err
	}
	npkt.WriteBody(rbyte)

	rb, err := codec.DefaultCodec.Marshal(rhdr)
	if err != nil {
		fmt.Println("Marshal(rsp)", err)
		return err
	}

	npkt.WriteBody(rb)

	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	// send to dispather
	return nil
}