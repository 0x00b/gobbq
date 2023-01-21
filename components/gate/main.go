package main

import (
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/server"
)

func main() {
	svr := gobbq.NewSever(server.WithPacketHandler(NewGatePacketHandler()))

	// RegisterTestEntity(svr, &TestEntity{})

	go svr.ListenAndServe(server.TCP, ":1234")
	go svr.ListenAndServe(server.KCP, ":1235")
	err := svr.ListenAndServe(server.WebSocket, ":8080")

	fmt.Println(err)
}