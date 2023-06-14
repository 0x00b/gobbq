package main

import (
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func (p *Proxy) HandleEOF(prw *nets.PacketReadWriter) {
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

func (p *Proxy) HandleTimeOut(prw *nets.PacketReadWriter) {
	p.HandleEOF(prw)
}

func (p *Proxy) HandleFail(prw *nets.PacketReadWriter) {
	p.HandleEOF(prw)
}
