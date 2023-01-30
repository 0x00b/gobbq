package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/conf"
)

func main() {

	fmt.Println(conf.C)

	// var te TestServerInterface = &TestEntity{}

	// RegisterTestService(te)

	// RegisterTestEntity(te)

	// svr := gobbq.NewSever(nets.WithPacketHandler(game.NewGamePacketHandler()))
	// go svr.ListenAndServe(nets.TCP, ":1234")
	// go svr.ListenAndServe(nets.KCP, ":1235")
	// err := svr.ListenAndServe(nets.WebSocket, ":8080")
	game.ConnectProxy()

	bufio.NewReader(os.Stdin).ReadString('\n')
	// fmt.Println(err)
}
