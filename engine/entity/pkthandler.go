package entity

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
)

var MethodNameError = errors.New("mothod name error")
var MethodNotFound = errors.New("mothod not found")
var ServiceNotFound = errors.New("service not found")
var EntityNotFound = errors.New("entity not found")
var UnknownCallType = errors.New("unknown call type")
var EmptyEntityID = errors.New("bad call, empty dst entity")

func NotMyMethod(err error) bool {
	return errors.Is(err, ServiceNotFound) || errors.Is(err, EntityNotFound)
}

var _ nets.PacketHandler = &MethodPacketHandler{}

type MethodPacketHandler struct {
}

func NewMethodPacketHandler() *MethodPacketHandler {
	st := &MethodPacketHandler{}
	return st
}

func (st *MethodPacketHandler) HandlePacket(pkt *codec.Packet) error {
	switch pkt.Header.ServiceType {
	case bbq.ServiceType_Entity:
		return st.handleCallEntity(pkt)
	case bbq.ServiceType_Service:
		return st.handleCallService(pkt)
	default:
	}

	return nil
}

func (st *MethodPacketHandler) handleCallService(pkt *codec.Packet) error {

	hdr := pkt.Header
	service := hdr.ServiceName

	svc, ok := Manager.Services[TypeName(service)]
	if !ok {
		return ServiceNotFound
	}

	svc.dispatchPkt(pkt)

	return nil
}

func (st *MethodPacketHandler) handleCallEntity(pkt *codec.Packet) error {

	eid := pkt.Header.GetDstEntity()
	if eid == "" {
		return EmptyEntityID
	}

	Manager.mu.RLock()
	defer Manager.mu.RUnlock()
	entity, ok := Manager.Entities[(EntityID(eid))]
	if !ok {
		return EntityNotFound
	}

	entity.dispatchPkt(pkt)

	return nil
}

// =========

// func HandleCallLocalMethod(pkt *codec.Packet, in any, callback func(c *Context, rsp any)) error {
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

// func handleLocalCallService(hdr *bbq.Header, in any, callback func(c *Context, rsp any)) error {
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
// 		return ServiceNotFound
// 	}

// 	ss.dispatchPkt(pkt)

// 	return nil
// 	// return handleCallMethod(nil, hdr, in, callback, ed.Desc())

// }

// func handleLocalCallEntity(hdr *bbq.Header, in any, callback func(c *Context, rsp any)) error {

// 	ety := hdr.GetDstEntity()
// 	if ety == nil {
// 		return EmptyEntityID
// 	}

// 	Manager.mu.RLock()
// 	defer Manager.mu.RUnlock()
// 	entity, ok := Manager.Entities[(EntityID(ety.ID))]
// 	if !ok {
// 		return EntityNotFound
// 	}
// 	entity.dispatchPkt(pkt)
// 	return nil

// 	// return handleCallMethod(entity.Context(), hdr, in, callback, entity.Desc())
// }

// func handleCallMethod(c *Context, hdr *bbq.Header, in any, callback func(c *Context, rsp any), sd *ServiceDesc) error {

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
