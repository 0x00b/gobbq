package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	xlog.Init("trace", true, true, os.Stdout, xlog.DefaultLogTag{})

	ex.ConnProxy(nets.WithPacketHandler(NewProxyPacketHandler()))
	entity.ProxyRegister = &RegisterProxy{}

	svr := gobbq.NewSever(nets.WithPacketHandler(NewClientPacketHandler()))

	gatepb.RegisterGateService(&GateService{})

	err := svr.ListenAndServe(
		nets.NetWorkName(conf.C.Gate.Inst[0].Net),
		fmt.Sprintf(":%s", conf.C.Gate.Inst[0].Port))

	fmt.Println(err)
}
