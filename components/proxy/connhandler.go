package main

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

type ConnHandler struct {
}

func (ch *ConnHandler) HandleEOF(prw *codec.PacketReadWriter) {
	func() {
		proxyInst.etyMtx.Lock()
		defer proxyInst.etyMtx.Unlock()
		for eid, v := range proxyInst.entityMaps {
			if v == prw {
				xlog.Debugln("remove eid:", eid)
				delete(proxyInst.entityMaps, eid)
			}
		}

	}()

	func() {
		proxyInst.svcMtx.Lock()
		defer proxyInst.svcMtx.Unlock()
		for svc, svcPrw := range proxyInst.svcMaps {
			if prw == svcPrw {
				delete(proxyInst.svcMaps, svc)
			}
		}
	}()
}

func (ch *ConnHandler) HandleTimeOut(prw *codec.PacketReadWriter) {
	ch.HandleEOF(prw)
}

func (ch *ConnHandler) HandleFail(prw *codec.PacketReadWriter) {

}
