package main

import (
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

type ProxyConnCallBack struct {
	gate *Gate
}

func (p *ProxyConnCallBack) HandleClose(cn *nets.Conn) {
	xlog.Info("proxy disconnect")
}

func (p *ProxyConnCallBack) HandleEOF(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *ProxyConnCallBack) HandleTimeOut(cn *nets.Conn) {
	p.HandleClose(cn)
}

func (p *ProxyConnCallBack) HandleFail(cn *nets.Conn) {
	p.HandleClose(cn)
}
