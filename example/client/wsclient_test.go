package main

import (
	"os"
	"testing"

	"github.com/0x00b/gobbq/components/gate/client"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/xlog"
)

func TestWSClient(m *testing.T) {

	xlog.Init("trace", true, true, os.Stdout, xlog.DefaultLogTag{})

	client := client.NewClient()

	es := exampb.NewEchoServiceClient(client.Gate)
	rsp, err := es.SayHello(client.Context(), &exampb.SayHelloRequest{Text: "hello"})
	if err != nil {
		panic(err)
	}

	xlog.Println("recv:", rsp)
}
