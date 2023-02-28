package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ProxySvc struct {
	*Proxy

	entity.Service
}

func main() {

	fmt.Println(conf.C)

	xlog.Init("info", false, true, &lumberjack.Logger{
		Filename:  "./proxy.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	p := NewProxy()

	proxypb.RegisterProxyEtyEntity(p.EntityMgr, p)
	proxypb.RegisterProxySvcService(p.EntityMgr, &ProxySvc{Proxy: p})

	p.RegisterNetService(nets.NewNetService(
		nets.WithPacketHandler(p),
		// nets.WithConnErrHandler(p),
		nets.WithNetwork(nets.TCP, fmt.Sprintf(":%s", conf.C.Proxy.Inst[0].Port))),
	)

	err := p.Serve()

	fmt.Println(err)
}
