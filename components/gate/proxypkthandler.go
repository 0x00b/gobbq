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

		rw, ok := func() (*codec.PacketReadWriter, bool) {
			gt.cltMtx.Lock()
			defer gt.cltMtx.Unlock()
			prw, ok := gt.cltMap[id.ID]
			return prw, ok
		}()
		if !ok {
			return errors.New("unknown client")
		}

		// todo 需要处理一下kcp的断开连接，否则会阻塞在这里，以及read也会阻塞，导致goroutine得不到释放
		// https://github.com/skywind3000/kcp/issues/176
		return rw.SendPackt(pkt)
	}

	return err

}
