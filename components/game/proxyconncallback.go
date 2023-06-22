package game

import (
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func (g *Game) HandleClose(cn *nets.Conn) {
	// proxy
	xlog.Info("proxy disconnect")
}

func (g *Game) HandleEOF(cn *nets.Conn) {
	g.HandleClose(cn)
}

func (g *Game) HandleTimeOut(cn *nets.Conn) {
	g.HandleClose(cn)
}

func (g *Game) HandleFail(cn *nets.Conn) {
	g.HandleClose(cn)
}
