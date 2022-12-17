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

	ws := codec.NewMessageReadWriter(context.Background(), wsc)

	pkt := codec.NewMessage()
	pkt.WriteBytes([]byte("test"))
	ws.WriteMessage(pkt)

	if pkt, err = ws.ReadMessage(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", string(pkt.MessageBody()))
}
