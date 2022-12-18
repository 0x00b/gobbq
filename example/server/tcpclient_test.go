package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
)

func TestTcpClient(m *testing.T) {

	wsc, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	ws := codec.NewPacketReadWriter(context.Background(), wsc)

	msg := codec.NewPacket()
	msg.WriteBytes([]byte("dsfsdfs"))
	ws.WritePacket(msg)

	if msg, err = ws.ReadPacket(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", string(msg.PacketBody()))
}
