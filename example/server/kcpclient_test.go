package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server/transport"
	"github.com/xtaci/kcp-go"
)

func TestKcpClient(m *testing.T) {

	wsc, err := kcp.DialWithOptions("127.0.0.1:1235", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("runing")

	ctx := context.Background()

	ct := transport.NewClientTransport(ctx, wsc)

	pkt := codec.NewMessage()
	pkt.WriteBytes([]byte("dsfsdfs"))

	fmt.Println("writing")
	ct.WriteMessage(pkt)
	fmt.Println("writed")

	ct.Serve()
}
