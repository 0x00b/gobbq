package main

import (
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func (p *Proxy) HandleClose(cn *nets.Conn) {

	// game or gate disconnect

	func() {
		p.instMtx.Lock()
		defer p.instMtx.Unlock()
		for eid, v := range p.instMaps {
			if v == cn {
				xlog.Debugln("remove eid:", eid)
				delete(p.instMaps, eid)
			}
		}
	}()

	func() {
		p.svcMtx.Lock()
		defer p.svcMtx.Unlock()
		for svc, svcPrw := range p.svcMaps {
			if cn == svcPrw {
				xlog.Debugln("remove svr:", svc)
				delete(p.svcMaps, svc)
			}
		}
	}()

}

func (p *Proxy) HandleEOF(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *Proxy) HandleTimeOut(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *Proxy) HandleFail(cn *nets.Conn) {
	p.HandleClose(cn)
}
