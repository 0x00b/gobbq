package main

import (
	"github.com/0x00b/gobbq/engine/codec"
)

type ConnHandler struct {
}

func (ch *ConnHandler) HandleEOF(prw *codec.PacketReadWriter) {
	etyMtx.Lock()
	defer etyMtx.Unlock()
	for eid, v := range entityMaps {
		if v == prw {
			// xlog.Println("remove eid:", eid)
			delete(entityMaps, eid)
		}
	}

	for svc, prws := range svcMap {
		idx := 0
		for i, t := range prws {
			if t != prw {
				if i != idx {
					// xlog.Println("remove svc:", prws[idx])
					prws[idx] = t
				}
				idx++
			}
		}
		svcMap[svc] = prws[:idx]
	}
}

func (ch *ConnHandler) HandleTimeOut(prw *codec.PacketReadWriter) {
	ch.HandleEOF(prw)
}

func (ch *ConnHandler) HandleFail(prw *codec.PacketReadWriter) {

}
