package nets_test

import (
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/xlog"
	"github.com/xtaci/kcp-go"
)

func TestWSServer(m *testing.T) {

	listener, err := kcp.ListenWithOptions("127.0.0.1:8899", nil, 10, 3)
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
				_, err := conn.Read(b[:])
				if err != nil {
					panic(err)
				}
				fmt.Println("recv", string(b[:]))
				t := "send:" + string(b[:])
				conn.Write([]byte(t))

				fmt.Println("send", string(b[:]))
			}

		}()
	}
}
