package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	conf.Init("proxy.yaml")

	fmt.Println(conf.C)

	proxyInst.ConnOtherProxy(nets.WithPacketHandler(NewProxyPacketHandler()))

	xlog.Init("trace", false, true, os.Stdout, xlog.DefaultLogTag{})

	svr := gobbq.NewSever(nets.WithPacketHandler(NewProxyPacketHandler()), nets.WithConnHandler(&ConnHandler{}))

	proxypb.RegisterProxyService(&ProxyService{})

	err := svr.ListenAndServe(nets.TCP, fmt.Sprintf(":%s", conf.C.Proxy.Inst[0].Port))

	fmt.Println(err)
}
