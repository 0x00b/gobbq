package main

import (
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	xlog.Init("trace", true, true, os.Stdout, xlog.DefaultLogTag{})

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

	xlog.Println("service", entity.GetPacket(c).Header.String(), req.String())

	echoClient := exampb.NewEchoEtyEntity(c)
	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", entity.GetPacket(c).Header.String(), rsp.String())

	return rsp, nil
}

type EchoEntity struct {
	entity.Entity
}

func (*EchoEntity) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("entity req", entity.GetPacket(c).Header.String(), req.String())

	client := exampb.NewClientEntityClient(req.CLientID)
	rsp, err := client.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("entity response:", entity.GetPacket(c).Header.String(), rsp.String())

	return rsp, nil
}
