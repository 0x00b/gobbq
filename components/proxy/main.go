package main

import (
	"fmt"
	"os"

	bbq "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	fmt.Println(conf.C)

	proxyInst.ConnOtherProxy(nets.WithPacketHandler(NewProxyPacketHandler()))

	xlog.Init("info", false, true, os.Stdout, xlog.DefaultLogTag{})

	proxypb.RegisterProxySvcService(&ProxyService{})

	svr := bbq.NewServer()

	svr.RegisterNetService(nets.NewNetService(
		nets.WithPacketHandler(NewProxyPacketHandler()),
		nets.WithConnHandler(&ConnHandler{}),
		nets.WithNetwork(nets.TCP),
		nets.WithAddress(fmt.Sprintf(":%s", conf.C.Proxy.Inst[0].Port))),
	)

	err := svr.ListenAndServe()

	fmt.Println(err)
}
