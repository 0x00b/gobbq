package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

type ProxySvc struct {
	*Proxy

	entity.Service
}

func main() {

	fmt.Println(conf.C)

	xlog.Init("trace", false, true, os.Stdout, xlog.DefaultLogTag{})

	p := NewProxy()

	proxypb.RegisterProxyEtyEntity(p.Server.EntityMgr, p)
	proxypb.RegisterProxySvcService(p.Server.EntityMgr, &ProxySvc{Proxy: p})

	p.Server.RegisterNetService(nets.NewNetService(
		nets.WithPacketHandler(p),
		// nets.WithConnErrHandler(p),
		nets.WithNetwork(nets.TCP, fmt.Sprintf(":%s", conf.C.Proxy.Inst[0].Port))),
	)

	err := p.Server.ListenAndServe()

	fmt.Println(err)
}
