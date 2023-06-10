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
		Filename:  "./server.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	fmt.Println(conf.C)

	var g = game.NewGame()

	// exampb.RegisterEchoService(&EchoService{})
	exampb.RegisterEchoService(g.EntityMgr, &EchoService{})

	exampb.RegisterEchoEtyEntity(g.EntityMgr, &EchoEntity{})

	g.Serve()
}

type EchoService struct {
	entity.Service
}

func (*EchoService) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("service", c.Packet().Header.String(), req.String())

	echoClient := exampb.NewEchoEtyEntity(c)
	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}

type EchoEntity struct {
	entity.Entity
}

func (e *EchoEntity) OnTick() {
	// xlog.Infoln("tick...")
}

func (e *EchoEntity) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	// e.AddCallback(1*time.Second, func() {
	// 	xlog.Infoln("tick.1111..")
	// })
	// e.AddTimer(2*time.Second, func() {
	// 	xlog.Infoln("tick..2222.")
	// })

	xlog.Println("entity req", c.Packet().Header.String(), req.String())

	client := exampb.NewClientEntityClient(req.CLientID)
	rsp, err := client.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}
