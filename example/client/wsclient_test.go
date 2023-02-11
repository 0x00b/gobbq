package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var wg sync.WaitGroup

type GamePacketHandler struct {
	entity.MethodPacketHandler
}

func TestWSClient(m *testing.T) {
	xlog.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./proxy.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})

	cfg := conf.C.Gate.Inst[0]
	client, err := nets.Connect(
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(&GamePacketHandler{}))
	if err != nil {
		panic(err)
	}

	// wg.Add(1)

	// // gate := gatepb.NewGateServiceClient(client)
	// // gate.RegisterClient(nil, &gatepb.RegisterClientRequest{}, func(c *entity.Context, rsp *gatepb.RegisterClientResponse) {
	// // 	fmt.Println("recv", rsp.String())

	// // 	es := exampb.NewEchoServiceClient(client)
	// // 	wg.Add(1)
	// // 	es.SayHello(c, &exampb.SayHelloRequest{Text: "hello"}, func(c *entity.Context, rsp *exampb.SayHelloResponse) {
	// // 		fmt.Println("recv", string(c.Packet().PacketBody()))
	// // 		fmt.Println(rsp)
	// // 		wg.Done()
	// // 	})

	// // })

	pkt, release := codec.NewPacket()
	defer release()

	hdr := pkt.Header

	hdr.Version = 1
	hdr.RequestId = snowflake.GenUUID()
	hdr.Timeout = 1
	hdr.RequestType = bbq.RequestType_RequestRequest
	hdr.ServiceType = bbq.ServiceType_Service
	hdr.SrcEntity = "eid"
	hdr.DstEntity = ""
	hdr.ServiceName = "exampb.EchoService"
	hdr.Method = "SayHello"
	hdr.ContentType = bbq.ContentType_Proto
	hdr.CompressType = bbq.CompressType_None
	hdr.CheckFlags = 0
	hdr.TransInfo = map[string][]byte{}
	hdr.ErrCode = 0
	hdr.ErrMsg = ""

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(&exampb.SayHelloRequest{
		Text: "test111",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("raw data len:", len(pkt.Data()), len(hdrBytes))
	pkt.WriteBody(hdrBytes)
	client.WritePacket(pkt)

	wg.Wait()
}
