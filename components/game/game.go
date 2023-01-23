package game

import (
	"context"
	"log"

	"github.com/0x00b/gobbq/engine/nets"
	"golang.org/x/net/websocket"
)

func test() {

	origin := "http://localhost:8080/"
	url := "ws://localhost:8080/ws"
	wsc, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	opts := nets.WithPacketHandler(NewGamePacketHandler())

	nets.NewClient(context.Background(), wsc, opts)
}
