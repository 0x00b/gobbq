package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
)

func main() {

	xlog.Init("trace", true, true, os.Stdout, xlog.DefaultLogTag{})

	fmt.Println(conf.C)

	ex.ConnProxy(nets.WithPacketHandler(game.NewGamePacketHandler()))
	entity.ProxyRegister = &game.RegisterProxy{}

	// exampb.RegisterEchoService(&EchoService{})
	exampb.RegisterEchoService(&EchoService{})

	exampb.RegisterEchoEtyEntity(&EchoEntity{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	bufio.NewReader(os.Stdin).ReadString('\n')
	// fmt.Println(err)
}

type EchoService struct {
	entity.Service
}

func (*EchoService) SayHello(c *entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("service", c.Packet().Header.String(), req.String())

	echoClient := exampb.NewEchoEtyEntity(c, ex.ProxyClient)
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

func (*EchoEntity) SayHello(c *entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("entity req", c.Packet().Header.String(), req.String())

	return &exampb.SayHelloResponse{Text: "echo entity response"}, nil
}
