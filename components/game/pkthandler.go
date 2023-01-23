package game

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto"
)

var _ nets.PacketHandler = &GamePacketHandler{}

type GamePacketHandler struct {
}

func NewGamePacketHandler() *GamePacketHandler {
	st := &GamePacketHandler{}
	// st.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, st)
	return st
}

func (st *GamePacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {
	switch pkt.GetHeader().CallType {
	case proto.CallType_CallEntity:
		return st.HandleEntity(c, pkt)
	case proto.CallType_CallService:
		return st.HandleService(c, pkt)
	default:
	}
	return errors.New("unknown call type")
}

func (st *GamePacketHandler) HandleMethod(c context.Context, pkt *codec.Packet, sd *entity.ServiceDesc) error {

	hdr := pkt.GetHeader()

	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		fmt.Println("err mothod")
		return errors.New("err mothod")
	}

	// service := sm[:pos]
	method := sm[pos+1:]

	mt := sd.Methods[method]
	dec := func(v interface{}) error {
		reqbuf := pkt.PacketBody()
		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
		return err
	}

	rsp, err := mt.Handler(sd.ServiceImpl, c, dec, nil)

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
		return err
	}
	return nil
}

func (st *GamePacketHandler) HandleService(c context.Context, pkt *codec.Packet) error {

	hdr := pkt.GetHeader()

	fmt.Println("recv RequestHeader:", hdr.String())

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

	ed, ok := entity.Manager.Services[entity.ServiceType(service)]
	if !ok {
		return errors.New("unknown service type")
	}

	return st.HandleMethod(c, pkt, ed)

}

func (st *GamePacketHandler) HandleEntity(c context.Context, pkt *codec.Packet) error {

	hdr := pkt.GetHeader()
	ety := hdr.GetDstEntity()
	if ety == nil {
		return errors.New("bad call, empty dst entity")
	}

	fmt.Println("recv RequestHeader:", hdr.String())

	sd, ok := entity.Manager.Entities[(entity.EntityID(ety.ID))]
	if !ok {
		return errors.New("unknown entity id")
	}

	return st.HandleMethod(c, pkt, sd)
}
