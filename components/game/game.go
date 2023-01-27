package game

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/engine/nets"
)

func ConnectProxy() {
	ex.ConnProxy(nets.WithPacketHandler(NewGamePacketHandler()))
}
