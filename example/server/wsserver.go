package main

import (
	"time"

	"github.com/0x00b/gobbq/components/game"
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

	var g = game.NewGame()

	xlog.Info("RegisterEchoService")
	exampb.RegisterEchoService(g.EntityMgr, &EchoService{})

	xlog.Info("RegisterEchoEtyEntity")
	exampb.RegisterEchoEtyEntity(g.EntityMgr, &EchoEntity{})
	xlog.Info("start")

	g.Serve()
}

type EchoService struct {
	entity.Service
}

func (e *EchoService) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("service", c.Packet().Header.String(), req.String())

	echoClient, err := exampb.NewEchoEty(c)
	if err != nil {
		return nil, err
	}

	xlog.Println("watch starting")
	e.Watch(echoClient.EntityID)
	xlog.Println("watch done")

	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("Service response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}

type EchoEntity struct {
	entity.Entity

	exampb.EchoPropertyModel
}

func (e *EchoEntity) OnInit() {
	// var db db.IDatabase

	// e.Text = "xxx" // _id

	// e.ModelInit(e.Context(), db)

	// e.SetText("xxxx")

	_ = e.Text
}

func (e *EchoEntity) OnTick() {
	// xlog.Infoln("tick...")
}

func (e *EchoEntity) OnDestroy() {

	// e.Destroy(e.Context())

	// xlog.Infoln("tick...")
	// e.Save()
}

func (e *EchoEntity) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	e.AddCallback(2*time.Second, func() {
		xlog.Infoln("tick.1111..")
		e.Stop()
	})
	// e.AddTimer(2*time.Second, func() {
	// 	xlog.Infoln("tick..2222.")
	// })

	xlog.Println("entity req", c.Packet().Header.String(), req.String())

	client := exampb.NewClientClient(entity.EntityID(req.GetCLientID()))
	rsp, err := client.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}
