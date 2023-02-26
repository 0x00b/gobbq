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

	g := game.NewGame()

	exampb.RegisterEchoSvc2Service(g.EntityMgr, &EchoService2{})

	g.Serve()
}

type EchoService2 struct {
	entity.Service
}

func (*EchoService2) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("service2222 req", c.Packet().Header.String(), req.String())

	echoClient := exampb.NewEchoServiceClient(g.EntityMgr, ex.ProxyClient.GetPacketReadWriter())
	rsp, err := echoClient.SayHello(c, req)
	if err != nil {
		return nil, err
	}
	xlog.Println("echo svc response:", c.Packet().Header.String(), rsp.String())

	return rsp, nil
}
