package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
	"github.com/0x00b/gobbq/proto"
)

var _ server.PacketHandler = &ServerPacketHandler{}

type ServerPacketHandler struct {
	opts *server.ServerOptions
}

func NewServerPacketHandler(ctx context.Context, conn net.Conn, opts *server.ServerOptions) *ServerPacketHandler {
	st := &ServerPacketHandler{opts: opts}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *ServerPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	switch pkt.GetHeader().CallType {
	case proto.CallType_entity:
	case proto.CallType_service:
	default:
	}

	// hdr := &proto.Header{}

	hdr := pkt.GetHeader()

	// codec.DefaultCodec.Unmarshal(pkt.PacketBody()[:pkt.GetMsgHeaderLen()], hdr)

	// fmt.Println("recv RequestHeader:", hdr.String())
	// fmt.Println("recv len:", pkt.GetMsgHeaderLen(), pkt.GetPacketBodyLen())
	// fmt.Println("recv data:", string(pkt.PacketBody()[pkt.GetMsgHeaderLen():pkt.GetPacketBodyLen()]))

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
	ed := st.opts.Entities[service]
	mt := ed.Methods[method]
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}

	rsp, err := mt.Handler(ed.ServiceImpl, c, dec, nil)

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

	rb, err := codec.DefaultCodec.Marshal(rsp)
	if err != nil {
		fmt.Println("Marshal(rsp)", err)
		return err
	}

	npkt.WriteBody(rb)

	err = pkt.Src.WritePacket(npkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}
