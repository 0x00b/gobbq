package entity

import (
	"errors"
	"unsafe"

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
var ErrBadRequest = errors.New("bad call, nil parameters")

func NotMyMethod(err error) bool {
	return errors.Is(err, ErrServiceNotFound) || errors.Is(err, ErrEntityNotFound)
}

func (st *EntityManager) HandlePacket(pkt *nets.Packet) error {

	if pkt.Header.ServiceType == bbq.ServiceType_Entity ||
		// 因为虽然service本地也有,但是如果entityid不是本地,说明是其他地方的同名service调用过来的,需要回到原地
		pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		return st.handleCallEntity(pkt)
	}

	return st.handleCallService(pkt)
}

func (st *EntityManager) handleCallService(pkt *nets.Packet) error {

	hdr := pkt.Header
	service := hdr.Type

	svc, ok := st.GetService(service)
	if !ok {
		xlog.Traceln("service not found in local", unsafe.Pointer(st), service)
		return ErrServiceNotFound
	}

	DispatchPkt(svc, pkt)

	return nil
}

func (st *EntityManager) handleCallEntity(pkt *nets.Packet) error {

	eid := DstEntity(pkt)
	if eid.Invalid() {
		xlog.Traceln("recv:", pkt.Header.RequestId, ErrEmptyEntityID)
		return ErrEmptyEntityID
	}

	// xlog.Traceln("start find entity")

	entity, ok := st.GetEntity(eid)
	if !ok {
		xlog.Traceln("entity not found in local", unsafe.Pointer(st), eid)
		return ErrEntityNotFound
	}

	// xlog.Traceln("dispatchPkt send:", pkt.Header.RequestId)

	DispatchPkt(entity, pkt)

	return nil
}

// =========

func (st *EntityManager) LocalCall(pkt *nets.Packet, in any, respChan chan any) error {

	hdr := pkt.Header
	if hdr.ServiceType == bbq.ServiceType_Entity ||
		// 因为虽然service本地也有,但是如果entityid不是本地,说明是其他地方的同名service调用过来的,需要回到原地
		hdr.RequestType == bbq.RequestType_RequestRespone {
		return st.handleLocalCallEntity(pkt, in, respChan)
	}

	return st.handleLocalCallService(pkt, in, respChan)

}

func (st *EntityManager) handleLocalCallService(pkt *nets.Packet, in any, respChan chan any) error {

	service := pkt.Header.Type

	ss, ok := st.GetService(service)
	if !ok {
		xlog.Traceln("service not found in local", unsafe.Pointer(st), service)
		return ErrServiceNotFound
	}

	xlog.Traceln("handleLocalCallService", pkt.Header.String())

	return ss.dispatchLocalCall(pkt, in, respChan)
}

func (st *EntityManager) handleLocalCallEntity(pkt *nets.Packet, in any, respChan chan any) error {

	// xlog.Traceln("handleLocalCallEntity 1", pkt.Header.String())

	eid := DstEntity(pkt)
	if eid.Invalid() {
		return ErrEmptyEntityID
	}

	entity, ok := st.GetEntity(eid)
	if !ok {
		xlog.Traceln("entity not found in local", unsafe.Pointer(st), eid.ID())
		return ErrEntityNotFound
	}
	// xlog.Traceln("handleLocalCallEntity 2", pkt.Header.String())

	return entity.dispatchLocalCall(pkt, in, respChan)
}
