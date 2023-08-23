package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	InitConfig()

	xlog.Init("trace", false, true, &lumberjack.Logger{
		Filename:  "./proxy.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	p := NewProxy()

	proxypb.RegisterProxySvcService(p.EntityMgr, &ProxySvc{proxy: p})

	p.RegisterNetService(nets.NewNetService(
		nets.WithPacketHandler(p),
		nets.WithConnCallback(p),
		nets.WithNetwork(nets.TCP, fmt.Sprintf(":%s", CFG.Port))),
	)

	err := p.Serve()

	fmt.Println(err)
}
