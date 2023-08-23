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
		// p.svcMtx.Lock()
		// defer p.svcMtx.Unlock()
		for svc, svcPrws := range p.svcMaps {
			for i, t := range svcPrws {
				if cn == t {
					xlog.Debugln("remove svr:", svc)
					lastIdx := len(svcPrws) - 1
					svcPrws[i] = svcPrws[lastIdx]
					p.svcMaps[svc] = svcPrws[:lastIdx]
				}
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
