package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"golang.org/x/net/websocket"
)

func TestWSClient(m *testing.T) {

	origin := "http://localhost/"
	url := "ws://localhost/ws"
	wsc, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	ws := codec.NewPacketReadWriter(context.Background(), wsc)

	msg := codec.NewPacket()
	msg.WriteBytes([]byte("test"))
	ws.WritePacket(msg)

	if msg, err = ws.ReadPacket(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", string(msg.PacketBody()))
}
