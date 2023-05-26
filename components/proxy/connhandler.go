package main

import (
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

func (p *Proxy) HandleEOF(prw *codec.PacketReadWriter) {
	func() {
		p.instMtx.Lock()
		defer p.instMtx.Unlock()
		for eid, v := range p.instMaps {
			if v == prw {
				xlog.Debugln("remove eid:", eid)
				delete(p.instMaps, eid)
			}
		}

	}()

	func() {
		p.svcMtx.Lock()
		defer p.svcMtx.Unlock()
		for svc, svcPrw := range p.svcMaps {
			if prw == svcPrw {
				delete(p.svcMaps, svc)
			}
		}
	}()
}

func (p *Proxy) HandleTimeOut(prw *codec.PacketReadWriter) {
	p.HandleEOF(prw)
}

func (p *Proxy) HandleFail(prw *codec.PacketReadWriter) {
	p.HandleEOF(prw)
}
