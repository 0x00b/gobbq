package entity

import (
	"context"
	"errors"
	"strings"

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
	return errors.Is(err, ServiceNotFound) ||
		errors.Is(err, EntityNotFound) ||
		errors.Is(err, MethodNotFound)
}

var _ nets.PacketHandler = &MethodPacketHandler{}

type MethodPacketHandler struct {
}

func NewMethodPacketHandler() *MethodPacketHandler {
	st := &MethodPacketHandler{}
	return st
}

func (st *MethodPacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {
	switch pkt.GetHeader().ServiceType {
	case bbq.ServiceType_Entity:
		return st.handleCallEntity(c, pkt)
	case bbq.ServiceType_Service:
		return st.handleCallService(c, pkt)
	default:
	}
	return UnknownCallType
}

func (st *MethodPacketHandler) handleCallMethod(c context.Context, pkt *codec.Packet, sd *ServiceDesc) error {

	hdr := pkt.GetHeader()

	// todo method name repeat get
	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		return MethodNameError
	}

	// service := sm[:pos]
	method := sm[pos+1:]

	mt, ok := sd.Methods[method]
	if !ok {
		return MethodNotFound
	}

	mt.Handler(sd.ServiceImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

func (st *MethodPacketHandler) handleCallService(c context.Context, pkt *codec.Packet) error {

	hdr := pkt.GetHeader()

	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		return MethodNameError
	}

	service := sm[:pos]

	ed, ok := Manager.Services[TypeName(service)]
	if !ok {
		return ServiceNotFound
	}

	return st.handleCallMethod(c, pkt, ed)

}

func (st *MethodPacketHandler) handleCallEntity(c context.Context, pkt *codec.Packet) error {

	hdr := pkt.GetHeader()
	ety := hdr.GetDstEntity()
	if ety == nil {
		return EmptyEntityID
	}

	sd, ok := Manager.Entities[(EntityID(ety.ID))]
	if !ok {
		return EntityNotFound
	}

	return st.handleCallMethod(c, pkt, sd)
}

// =========

func HandleCallLocalMethod(c context.Context, hdr *bbq.Header, in interface{}, callback func(c context.Context, rsp interface{})) error {

	switch hdr.ServiceType {
	case bbq.ServiceType_Entity:
		return handleCallEntity(c, hdr, in, callback)
	case bbq.ServiceType_Service:
		return handleCallService(c, hdr, in, callback)
	default:
	}
	return UnknownCallType
}

func handleCallService(c context.Context, hdr *bbq.Header, in interface{}, callback func(c context.Context, rsp interface{})) error {
	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		return MethodNameError
	}

	service := sm[:pos]

	ed, ok := Manager.Services[TypeName(service)]
	if !ok {
		return ServiceNotFound
	}

	return handleCallMethod(c, hdr, in, callback, ed)

}

func handleCallEntity(c context.Context, hdr *bbq.Header, in interface{}, callback func(c context.Context, rsp interface{})) error {

	ety := hdr.GetDstEntity()
	if ety == nil {
		return EmptyEntityID
	}

	sd, ok := Manager.Entities[(EntityID(ety.ID))]
	if !ok {
		return EntityNotFound
	}

	return handleCallMethod(c, hdr, in, callback, sd)
}

func handleCallMethod(c context.Context, hdr *bbq.Header, in interface{}, callback func(c context.Context, rsp interface{}), sd *ServiceDesc) error {

	sm := hdr.GetMethod()
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		return MethodNameError
	}

	// service := sm[:pos]
	method := sm[pos+1:]

	mt, ok := sd.Methods[method]
	if !ok {
		return MethodNotFound
	}

	mt.LocalHandler(sd.ServiceImpl, c, in, callback, chainServerInterceptors(sd.interceptors))

	return nil
}
