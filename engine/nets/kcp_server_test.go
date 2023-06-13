package nets_test

import (
	"fmt"
	"os"
	"testing"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

type TestPacket struct {
}

func (tp *TestPacket) HandlePacket(pkt *codec.Packet) error {

	pkt.Src.SendPacket(pkt)

	fmt.Println(pkt.String())

	return nil
}

func TestKcpServer(m *testing.T) {

	xlog.Init("trace", true, true, os.Stdout)

	svr := bs.NewServer()
	svr.RegisterNetService(
		nets.NewNetService(
			nets.WithPacketHandler(&TestPacket{}),
			nets.WithNetwork(nets.KCP, ":8899")),
	)

	svr.ListenAndServe()

	// listener, err := kcp.ListenWithOptions("127.0.0.1:8899", nil, 10, 3)
	// if err != nil {
	// 	panic(err)
	// }
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		if err != nil {
	// 			continue
	// 		} else {
	// 			panic(err)
	// 		}
	// 	}

	// 	xlog.Infof("Connection from: %s", conn.RemoteAddr())
	// 	secure.GO( func() {
	// 		for {

	// 			var b [1024]byte
	// 			_, err := conn.Read(b[:])
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			req := new(exampb.SayHelloRequest)
	// 			codec.DefaultCodec.Unmarshal(b[:], req)

	// 			fmt.Println("recv", req.String())

	// 			t := "send:" + string(b[:])
	// 			conn.Write([]byte(t))

	// 			fmt.Println("send", string(b[:]))
	// 		}

	// 	})
	// }
}
