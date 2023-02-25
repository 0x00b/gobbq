package main

import (
	"fmt"
	"os"

	bbq "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	xlog.Init("debug", true, true, os.Stdout, xlog.DefaultLogTag{})

	ex.ConnProxy(nets.WithPacketHandler(NewProxyPacketHandler()))
	client := proxypb.NewProxySvcServiceClient(ex.ProxyClient.GetPacketReadWriter())

	rsp, err := client.RegisterInst(Inst.Context(), &proxypb.RegisterInstRequest{
		InstID: Inst.EntityID().ID,
	})
	if err != nil {
		panic(err)
	}

	Inst.ProxyID = rsp.ProxyID

	gatepb.RegisterGateService(&GateService{})

	svr := bbq.NewServer()

	svr.RegisterNetService(
		nets.NewNetService(
			nets.WithPacketHandler(NewClientPacketHandler()),
			nets.WithNetwork(nets.NetWorkName(conf.C.Gate.Inst[0].Net)),
			nets.WithAddress(fmt.Sprintf(":%s", conf.C.Gate.Inst[0].Port))),
	)

	err = svr.ListenAndServe()

	fmt.Println(err)
}
