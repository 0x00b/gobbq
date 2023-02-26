package entity

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
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

func (st *EntityManager) HandlePacket(pkt *codec.Packet) error {
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
	xlog.Traceln("recv response:", pkt.Header.RequestId)
	return st.handleCallEntity(pkt)
}

func (st *EntityManager) handleCallService(pkt *codec.Packet) error {

	hdr := pkt.Header
	service := hdr.DstEntity.Type

	svc, ok := st.Services[service]
	if !ok {
		return ErrServiceNotFound
	}

	svc.dispatchPkt(pkt)

	return nil
}

func (st *EntityManager) handleCallEntity(pkt *codec.Packet) error {

	eid := pkt.Header.GetDstEntity()
	if eid == nil && eid.ID == "" {
		xlog.Traceln("recv:", pkt.Header.RequestId, ErrEmptyEntityID)
		return ErrEmptyEntityID
	}

	xlog.Traceln("start find entity")

	st.mu.RLock()
	defer st.mu.RUnlock()
	entity, ok := st.Entities[eid.ID]
	if !ok {
		xlog.Traceln("recv no entity:", pkt.Header.RequestId, ErrEmptyEntityID)
		return ErrEntityNotFound
	}

	xlog.Traceln("dispatchPkt send:", pkt.Header.RequestId)

	entity.dispatchPkt(pkt)

	return nil
}

// =========

func (st *EntityManager) HandleCallLocalMethod(pkt *codec.Packet, in any, respChan chan any) error {
	switch pkt.Header.ServiceType {
	case bbq.ServiceType_Entity:
		return st.handleLocalCallEntity(pkt, in, respChan)
	case bbq.ServiceType_Service:
		return st.handleLocalCallService(pkt, in, respChan)
	default:
	}
	return ErrUnknownCallType
}

func (st *EntityManager) handleLocalCallService(pkt *codec.Packet, in any, respChan chan any) error {

	service := pkt.Header.DstEntity.Type

	ss, ok := st.Services[service]
	if !ok {
		return ErrServiceNotFound
	}

	xlog.Traceln("handleLocalCallService", pkt.Header.String())

	return ss.dispatchLocalCall(pkt, in, respChan)
}

func (st *EntityManager) handleLocalCallEntity(pkt *codec.Packet, in any, respChan chan any) error {

	xlog.Traceln("handleLocalCallEntity 1", pkt.Header.String())

	ety := pkt.Header.GetDstEntity()
	if ety == nil {
		return ErrEmptyEntityID
	}

	st.mu.RLock()
	defer st.mu.RUnlock()
	entity, ok := st.Entities[(ety.ID)]
	if !ok {
		return ErrEntityNotFound
	}
	xlog.Traceln("handleLocalCallEntity 2", pkt.Header.String())

	return entity.dispatchLocalCall(pkt, in, respChan)
}
