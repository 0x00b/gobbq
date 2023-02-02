package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/example/exampb"
)

func main() {

	fmt.Println(conf.C)

	ex.ConnProxy(nets.WithPacketHandler(game.NewGamePacketHandler()))
	entity.ProxyRegister = &game.RegisterProxy{}

	// exampb.RegisterEchoService(&EchoService{})
	exampb.RegisterEchoService(&EchoService{})

	exampb.RegisterEchoEtyEntity(&EchoEntity{})

	bufio.NewReader(os.Stdin).ReadString('\n')
	// fmt.Println(err)
}

type EchoService struct {
	entity.Service
}

func (*EchoService) SayHello(c *entity.Context, req *exampb.SayHelloRequest, ret func(*exampb.SayHelloResponse, error)) {

	fmt.Println("service", req.String())

	echoClient := exampb.NewEchoEtyEntity(c, ex.ProxyClient)
	echoClient.SayHello(c, req, func(c *entity.Context, rsp *exampb.SayHelloResponse) {
		ret(rsp, nil)
	})
}

type EchoEntity struct {
	entity.Entity
}

func (*EchoEntity) SayHello(c *entity.Context, req *exampb.SayHelloRequest, ret func(*exampb.SayHelloResponse, error)) {

	fmt.Println("entity req", req.String())

	ret(&exampb.SayHelloResponse{Text: "echo entity response"}, nil)
}
