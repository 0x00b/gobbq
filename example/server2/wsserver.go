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

	xlog.Init("trace", true, true, &lumberjack.Logger{
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

func (*EchoService2) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("service2222 req", entity.GetPacket(c).Header.String(), req.String())

	echoClient := exampb.NewEchoServiceClient()
	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("echo svc response:", entity.GetPacket(c).Header.String(), rsp.String())

	return rsp, nil
}
