package main

import (
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func (p *Gate) HandleClose(cn *nets.Conn) {
	p.cltMtx.Lock()
	defer p.cltMtx.Unlock()
	for eid, v := range p.cltMap {
		if v == cn {
			xlog.Debugln("remove client:", eid)
			delete(p.cltMap, eid)

			// do something
		}
	}
}

func (p *Gate) HandleEOF(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *Gate) HandleTimeOut(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *Gate) HandleFail(cn *nets.Conn) {
	p.HandleClose(cn)
}
