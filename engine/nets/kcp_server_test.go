package nets_test

import (
	"fmt"
	"testing"

	bs "github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
)

type TestPacket struct {
}

func (tp *TestPacket) HandlePacket(pkt *codec.Packet) error {

	fmt.Println(pkt.String())

	pkt.Src.SendPackt(pkt)

	return nil
}

func TestKcpServer(m *testing.T) {

	svr := bs.NewServer()
	svr.RegisterNetService(
		nets.NewNetService(
			nets.WithPacketHandler(&TestPacket{}),
			nets.WithNetwork(nets.NetWorkName("kcp"), fmt.Sprintf(":8899"))),
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
	// 	go func() {
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

	// 	}()
	// }
}
