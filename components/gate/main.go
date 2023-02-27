package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	xlog.Init("trace", true, true, os.Stdout, xlog.DefaultLogTag{})

	gt := NewGate()

	gatepb.RegisterGateService(gt.EntityMgr, gt)

	gt.RegisterNetService(
		nets.NewNetService(
			nets.WithPacketHandler(NewClientPacketHandler(gt.EntityMgr)),
			nets.WithNetwork(nets.NetWorkName(conf.C.Gate.Inst[0].Net), fmt.Sprintf(":%s", conf.C.Gate.Inst[0].Port))),
	)

	err := gt.Serve()

	fmt.Println(err)
}
