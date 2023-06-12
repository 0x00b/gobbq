package main

import (
	"fmt"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	xlog.Init("info", true, true, &lumberjack.Logger{
		Filename:  "./server2.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	fmt.Println(conf.C)

	var g = game.NewGame()

	exampb.RegisterEchoSvc2Service(g.EntityMgr, &EchoService2{})

	g.Serve()
}

type EchoService2 struct {
	entity.Service
}

func (e *EchoService2) OnTick() {
	// xlog.Infoln("tick...")
}

func (e *EchoService2) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	// e.AddCallback(1*time.Second, func() {
	// 	xlog.Infoln("tick.1111..")
	// })
	// e.AddTimer(2*time.Second, func() {
	// 	xlog.Infoln("tick..2222.")
	// })

	xlog.Println("service2222 req", c.Packet().Header.String(), req.String())

	echoClient := exampb.NewEchoServiceClient()
	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("echo svc response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}
