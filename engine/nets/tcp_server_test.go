package nets_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/0x00b/gobbq/xlog"
)

func TestTcpServer(m *testing.T) {

	listener, err := net.Listen("tcp", "127.0.0.1:8899")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err != nil {
				continue
			} else {
				panic(err)
			}
		}

		xlog.Infof("Connection from: %s", conn.RemoteAddr())
		go func() {
			for {

				var b [1024]byte
				len, err := conn.Read(b[:])
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("recv & send", len, string(b[:len]))
				conn.Write([]byte(b[:len]))

				fmt.Println("send", string(b[:]))
			}

		}()
	}
}
