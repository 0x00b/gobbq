package main

import (
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/nets"
)

func main() {
	svr := gobbq.NewSever(nets.WithPacketHandler(NewClientPacketHandler()))

	// RegisterTestEntity(svr, &TestEntity{})

	go svr.ListenAndServe(nets.TCP, ":1234")
	go svr.ListenAndServe(nets.KCP, ":1235")
	err := svr.ListenAndServe(nets.WebSocket, ":8080")

	fmt.Println(err)
}
