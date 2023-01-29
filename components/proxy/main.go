package main

import (
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
)

func main() {
	fmt.Println(conf.C)

	svr := gobbq.NewSever(nets.WithPacketHandler(NewProxyPacketHandler()))

	proxypb.RegisterProxyService(&ProxyService{})

	err := svr.ListenAndServe(nets.TCP, fmt.Sprintf(":%s", conf.C.Proxy.Inst[0].Port))

	fmt.Println(err)
}
