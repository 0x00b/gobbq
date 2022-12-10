package gobbq_test

import (
	"fmt"
	"testing"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/server"
)

func TestMain(m *testing.T) {
	// net.Listen("", "")
	// net.ListenPacket("", "")
	svr := gobbq.NewSever()
	err := svr.ListenAndServe(server.WebSocket, ":80")
	fmt.Println(err)

}
