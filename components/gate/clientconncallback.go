package main

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

// todo  kcp 协议链接断开感知
func (p *Gate) HandleClose(cn *nets.Conn) {
	p.cltMtx.Lock()
	defer p.cltMtx.Unlock()
	for eid, v := range p.cltMap {
		if v == cn {
			xlog.Debugln("remove client:", eid)
			delete(p.cltMap, eid)

			wc := p.watcher[eid]

			xlog.Traceln("watcher:", wc)

			for id := range wc {
				client := entity.NewBbqSysEntityClient(id)
				client.SysNotify(p.Context(), &entity.WatchRequest{EntityID: uint64(eid)})
			}
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
