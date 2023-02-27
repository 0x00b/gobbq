package main

import (
	"errors"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

var _ nets.PacketHandler = &Gate{}

func (gt *Gate) HandlePacket(pkt *codec.Packet) error {

	err := gt.Server.EntityMgr.HandlePacket(pkt)
	if err == nil {
		// handle succ
		return nil
	}

	if entity.NotMyMethod(err) {
		// send to client
		id := pkt.Header.GetDstEntity()

		gt.cltMtx.Lock()
		defer gt.cltMtx.Unlock()

		rw, ok := gt.cltMap[id.ID]
		if !ok {
			return errors.New("unknown client")
		}

		return rw.SendPackt(pkt)
	}

	return err

}
