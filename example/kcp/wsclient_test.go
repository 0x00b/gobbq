package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/xtaci/kcp-go"
	//这里使用的是 gorilla 的 websocket 库
)

func TestWSClient(m *testing.T) {

	wsc, err := kcp.DialWithOptions("127.0.0.1:1234", nil, 10, 3)
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
