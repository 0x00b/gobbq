package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/0x00b/gobbq/components/gate/gatepb"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type GateEntity struct {
	entity.Entity //没啥用
	gate          *Gate
}

// 使用service的entity id
func (g *GateEntity) EntityID() entity.EntityID {
	return g.gate.EntityID()
}

func main() {

	secure.GO(func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	})

	xlog.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./gate.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	gt := NewGate()

	gatepb.RegisterGateService(gt.EntityMgr, gt)

	gt.RegisterNetService(
		nets.NewNetService(
			nets.WithPacketHandler(NewClientPacketHandler(gt)),
			nets.WithConnCallback(gt),
			nets.WithNetwork(nets.NetWorkName(conf.C.Gate.Inst[0].Net), fmt.Sprintf(":%s", conf.C.Gate.Inst[0].Port))),
	)

	err := gt.Serve()

	fmt.Println(err)
}
