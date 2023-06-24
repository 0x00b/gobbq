package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &Gate{}

func (gt *Gate) isMyPacket(pkt *nets.Packet) bool {

	hdr := pkt.Header
	dstEty := entity.DstEntity(pkt)
	if hdr.GetServiceType() == bbq.ServiceType_Entity ||
		hdr.RequestType == bbq.RequestType_RequestRespone {
		_, ok := gt.EntityMgr.GetEntity(dstEty)
		return ok
	}

	_, ok := gt.EntityMgr.GetService(hdr.GetType())
	return ok

}

func (gt *Gate) HandlePacket(pkt *nets.Packet) error {

	if gt.isMyPacket(pkt) {
		err := gt.Server.EntityMgr.HandlePacket(pkt)
		if err != nil {
			xlog.Errorln("bad req handle:", pkt.Header.String(), err)
		}
		return err
	}

	dstEty := entity.DstEntity(pkt)
	rw, ok := gt.GetClient(dstEty)
	if !ok {
		return errors.New("unknown client")
	}

	// todo 需要处理一下kcp的断开连接，否则会阻塞在这里，以及read也会阻塞，导致goroutine得不到释放
	// https://github.com/skywind3000/kcp/issues/176
	return rw.SendPacket(pkt)

}

func (gt *Gate) GetClient(eid entity.EntityID) (*nets.Conn, bool) {
	gt.cltMtx.Lock()
	defer gt.cltMtx.Unlock()
	prw, ok := gt.cltMap[eid.ID()]
	return prw, ok
}
