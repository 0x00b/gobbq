package game

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &GamePacketHandler{}

type GamePacketHandler struct {
	*entity.MethodPacketHandler
}

func NewGamePacketHandler() *GamePacketHandler {
	st := &GamePacketHandler{entity.NewMethodPacketHandler()}
	return st
}

// func (st *GamePacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {
// 	switch pkt.GetHeader().ServiceType {
// 	case bbq.ServiceType_Entity:
// 		return st.HandleCallEntity(c, pkt)
// 	case bbq.ServiceType_Service:
// 		return st.HandleCallService(c, pkt)
// 	default:
// 	}
// 	return errors.New("unknown call type")
// }

// func (st *GamePacketHandler) HandleCallMethod(c context.Context, pkt *codec.Packet, sd *entity.ServiceDesc) error {

// 	hdr := pkt.GetHeader()

// 	sm := hdr.GetMethod()
// 	if sm != "" && sm[0] == '/' {
// 		sm = sm[1:]
// 	}
// 	pos := strings.LastIndex(sm, "/")
// 	if pos == -1 {
// 		fmt.Println("err mothod")
// 		return errors.New("err mothod")
// 	}

// 	// service := sm[:pos]
// 	method := sm[pos+1:]

// 	mt := sd.Methods[method]
// 	dec := func(v interface{}) error {
// 		reqbuf := pkt.PacketBody()
// 		err := codec.GetCodec(hdr.GetContentType()).Unmarshal(reqbuf, v)
// 		return err
// 	}

// 	mt.Handler(sd.ServiceImpl, c, dec, nil)

// 	return nil
// }

// func (st *GamePacketHandler) HandleCallService(c context.Context, pkt *codec.Packet) error {

// 	hdr := pkt.GetHeader()

// 	fmt.Println("recv RequestHeader:", hdr.String())

// 	sm := hdr.GetMethod()
// 	if sm != "" && sm[0] == '/' {
// 		sm = sm[1:]
// 	}
// 	pos := strings.LastIndex(sm, "/")
// 	if pos == -1 {
// 		fmt.Println("err mothod")
// 		return errors.New("err mothod")
// 	}

// 	service := sm[:pos]

// 	ed, ok := entity.Manager.Services[entity.TypeName(service)]
// 	if !ok {
// 		return errors.New("unknown service type")
// 	}

// 	return st.HandleCallMethod(c, pkt, ed)

// }

// func (st *GamePacketHandler) HandleCallEntity(c context.Context, pkt *codec.Packet) error {

// 	hdr := pkt.GetHeader()
// 	ety := hdr.GetDstEntity()
// 	if ety == nil {
// 		return errors.New("bad call, empty dst entity")
// 	}

// 	fmt.Println("recv RequestHeader:", hdr.String())

// 	sd, ok := entity.Manager.Entities[(entity.EntityID(ety.ID))]
// 	if !ok {
// 		return errors.New("unknown entity id")
// 	}

// 	return st.HandleCallMethod(c, pkt, sd)
// }
