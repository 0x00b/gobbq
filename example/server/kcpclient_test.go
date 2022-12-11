package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/xtaci/kcp-go"
)

func TestKcpClient(m *testing.T) {

	wsc, err := kcp.DialWithOptions("127.0.0.1:1235", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("runing")

	ws := codec.NewPacketReadWriter(context.Background(), wsc)

	pkt := codec.NewPacket()
	pkt.WriteBytes([]byte("dsfsdfs"))

	fmt.Println("writing")
	ws.WritePacket(pkt)
	fmt.Println("writed")

	fmt.Println("reading")
	if pkt, err = ws.ReadPacket(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("readed")
	fmt.Printf("Received: %s.\n", string(pkt.PacketBody()))
}
