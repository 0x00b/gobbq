package nets_test

import (
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
	"github.com/xtaci/kcp-go"
)

func TestWSServer(m *testing.T) {

	listener, err := kcp.ListenWithOptions("127.0.0.1:8899", nil, 0, 0)
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
		secure.GO(func() {
			for {

				var b [1024]byte
				n, err := conn.Read(b[:])
				if err != nil {
					panic(err)
				}
				fmt.Println("recv", string(b[:]))
				conn.Write([]byte(b[:n]))

				fmt.Println("send", string(b[:]))
			}

		})
	}
}

// func TestUdpServer(m *testing.T) {

// 	// listener, err := kcp.ListenWithOptions("127.0.0.1:8899", nil, 10, 3)

// 	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8899")
// 	listener, err := net.ListenUDP("udp", udpAddr)

// 	if err != nil {
// 		panic(err)
// 	}

// 	for {

// 		var b [1024]byte
// 		n, err := listener.Read(b[:])
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("recv", string(b[:]))
// 		listener.Write([]byte(b[:n]))

// 		fmt.Println("send", string(b[:]))
// 	}

// }
