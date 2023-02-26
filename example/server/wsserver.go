package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
)

var g = game.NewGame()

func main() {

	xlog.Init("info", true, true, os.Stdout, xlog.DefaultLogTag{})

	fmt.Println(conf.C)

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

	echoClient := exampb.NewEchoEtyEntity(c, g.EntityMgr, ex.ProxyClient.GetPacketReadWriter())
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

func (*EchoEntity) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("entity req", c.Packet().Header.String(), req.String())

	client := exampb.NewClientEntityClient(ex.ProxyClient.GetPacketReadWriter(), g.EntityMgr, req.CLientID)
	rsp, err := client.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}
