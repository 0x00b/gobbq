package entity

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var ErrMethodNameError = errors.New("mothod name error")
var ErrMethodNotFound = errors.New("mothod not found")
var ErrServiceNotFound = errors.New("service not found")
var ErrEntityNotFound = errors.New("entity not found")
var ErrUnknownCallType = errors.New("unknown call type")
var ErrEmptyEntityID = errors.New("bad call, empty dst entity")

func NotMyMethod(err error) bool {
	return errors.Is(err, ErrServiceNotFound) || errors.Is(err, ErrEntityNotFound)
}

var _ nets.PacketHandler = &MethodPacketHandler{}

type MethodPacketHandler struct {
}

func NewMethodPacketHandler() *MethodPacketHandler {
	st := &MethodPacketHandler{}
	return st
}

func (st *MethodPacketHandler) HandlePacket(pkt *codec.Packet) error {
	if pkt.Header.RequestType == bbq.RequestType_RequestRequest {
		switch pkt.Header.ServiceType {
		case bbq.ServiceType_Entity:
			return st.handleCallEntity(pkt)
		case bbq.ServiceType_Service:
			return st.handleCallService(pkt)
		default:
		}
		return nil
	}
	// response
	xlog.Println("recv response:", pkt.Header.RequestId)
	return st.handleCallEntity(pkt)
}

func (st *MethodPacketHandler) handleCallService(pkt *codec.Packet) error {

	hdr := pkt.Header
	service := hdr.DstEntity.Type

	svc, ok := Manager.Services[TypeName(service)]
	if !ok {
		return ErrServiceNotFound
	}

	svc.dispatchPkt(pkt)

	return nil
}

func (st *MethodPacketHandler) handleCallEntity(pkt *codec.Packet) error {

	eid := pkt.Header.GetDstEntity()
	if eid == nil && eid.ID == "" {
		xlog.Println("recv:", pkt.Header.RequestId, ErrEmptyEntityID)
		return ErrEmptyEntityID
	}

	xlog.Println("start find entity")

	Manager.mu.RLock()
	defer Manager.mu.RUnlock()
	entity, ok := Manager.Entities[eid.ID]
	if !ok {
		xlog.Println("recv no entity:", pkt.Header.RequestId, ErrEmptyEntityID)
		return ErrEntityNotFound
	}

	xlog.Println("dispatchPkt send:", pkt.Header.RequestId)

	entity.dispatchPkt(pkt)

	return nil
}

// =========

// func HandleCallLocalMethod(pkt *codec.Packet, in any, callback func(c Context, rsp any)) error {
// 	hdr := pkt.Header
// 	switch hdr.ServiceType {
// 	case bbq.ServiceType_Entity:
// 		return handleLocalCallEntity(hdr, in, callback)
// 	case bbq.ServiceType_Service:
// 		return handleLocalCallService(hdr, in, callback)
// 	default:
// 	}
// 	return UnknownCallType
// }

// func handleLocalCallService(hdr *bbq.Header, in any, callback func(c Context, rsp any)) error {
// 	sm := hdr.GetMethod()
// 	if sm != "" && sm[0] == '/' {
// 		sm = sm[1:]
// 	}
// 	pos := strings.LastIndex(sm, "/")
// 	if pos == -1 {
// 		return MethodNameError
// 	}

// 	service := sm[:pos]

// 	ss, ok := Manager.Services[TypeName(service)]
// 	if !ok {
// 		return ErrServiceNotFound
// 	}

// 	ss.dispatchPkt(pkt)

// 	return nil
// 	// return handleCallMethod(nil, hdr, in, callback, ed.Desc())

// }

// func handleLocalCallEntity(hdr *bbq.Header, in any, callback func(c Context, rsp any)) error {

// 	ety := hdr.GetDstEntity()
// 	if ety == nil {
// 		return EmptyEntityID
// 	}

// 	Manager.mu.RLock()
// 	defer Manager.mu.RUnlock()
// 	entity, ok := Manager.Entities[(EntityID(ety.ID))]
// 	if !ok {
// 		return ErrEntityNotFound
// 	}
// 	entity.dispatchPkt(pkt)
// 	return nil

// 	// return handleCallMethod(entity.Context(), hdr, in, callback, entity.Desc())
// }

// func handleCallMethod(c Context, hdr *bbq.Header, in any, callback func(c Context, rsp any), sd *ServiceDesc) error {

// 	sm := hdr.GetMethod()
// 	if sm != "" && sm[0] == '/' {
// 		sm = sm[1:]
// 	}
// 	pos := strings.LastIndex(sm, "/")
// 	if pos == -1 {
// 		return MethodNameError
// 	}

// 	// service := sm[:pos]
// 	method := sm[pos+1:]

// 	mt, ok := sd.Methods[method]
// 	if !ok {
// 		return MethodNotFound
// 	}

// 	mt.LocalHandler(sd.ServiceImpl, c, in, callback, chainServerInterceptors(sd.interceptors))

// 	return nil
// }
