package main

import (
	"os"
	"testing"

	"github.com/0x00b/gobbq/components/gate/client"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
)

type ClientService struct {
	entity.Entity
}

func (*ClientService) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("server req", c.Packet().Header.String(), req.String())

	return &exampb.SayHelloResponse{Text: "client service response"}, nil
}

func TestWSClient(m *testing.T) {

	xlog.Init("debug", true, true, os.Stdout, xlog.DefaultLogTag{})
	conf.Init("client.yaml")

	client := client.NewClient(&exampb.ClientEntityDesc, &ClientService{})

	es := exampb.NewEchoSvc2ServiceClient(client.Gate.GetPacketReadWriter())
	rsp, err := es.SayHello(client.Context(), &exampb.SayHelloRequest{
		Text:     "hello",
		CLientID: client.EntityID(),
	})
	if err != nil {
		panic(err)
	}

	xlog.Println("recv:", rsp)
}
