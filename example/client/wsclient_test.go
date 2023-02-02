package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/example/exampb"
)

var wg sync.WaitGroup

type GamePacketHandler struct {
	entity.MethodPacketHandler
}

func TestWSClient(m *testing.T) {

	cfg := conf.C.Gate.Inst[0]
	client, err := nets.Connect(
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(&GamePacketHandler{}))
	if err != nil {
		panic(err)
	}

	es := exampb.NewEchoServiceClient(client)
	wg.Add(1)
	es.SayHello(nil, &exampb.SayHelloRequest{Text: "hello"}, func(c *entity.Context, rsp *exampb.SayHelloResponse) {
		fmt.Println("recv", string(c.Packet().PacketBody()))
		fmt.Println(rsp)
		wg.Done()
	})

	// pkt, release := codec.NewPacket()
	// defer release()

	// hdr := pkt.Header

	// hdr.Version = 1
	// hdr.RequestId = "1"
	// hdr.Timeout = 1
	// hdr.RequestType = 0
	// hdr.ServiceType = 0
	// hdr.SrcEntity = &bbq.EntityID{ID: "222"}
	// hdr.DstEntity = &bbq.EntityID{ID: "111"}
	// hdr.Method = "new_client"
	// hdr.ContentType = 0
	// hdr.CompressType = 0
	// hdr.CheckFlags = codec.FlagDataChecksumIEEE
	// hdr.TransInfo = map[string][]byte{"xxx": []byte("22222")}
	// hdr.ErrCode = 0
	// hdr.ErrMsg = ""

	// pkt.WriteBody(nil)

	// wg.Add(1)
	// client.WritePacket(pkt)

	// hdr.Method = "helloworld.Test/SayHello"
	// hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(hdr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("raw data len:", len(pkt.Data()), len(hdrBytes))
	// pkt.WriteBody(hdrBytes)
	// client.WritePacket(pkt)

	wg.Wait()
}
