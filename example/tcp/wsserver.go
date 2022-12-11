package main

import (
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/server"
)

func main() {
	// net.Listen("", "")
	// net.ListenPacket("", "")
	svr := gobbq.NewSever()
	err := svr.ListenAndServe(server.TCP, ":1234")
	fmt.Println(err)

}
