package main

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto"
)

var wg sync.WaitGroup

type GamePacketHandler struct {
}

func (st *GamePacketHandler) HandlePacket(c context.Context, pkt *codec.Packet) error {

	hdr := &proto.Header{}

	err := codec.DefaultCodec.Unmarshal(pkt.Data(), hdr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("recv:", hdr.String())
	wg.Done()
	return nil
}

func TestWSClient(m *testing.T) {

	cfg := conf.C.Gate.Inst[0]
	client, err := nets.Connect(context.Background(),
		nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, nets.WithPacketHandler(&GamePacketHandler{}))

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		TransInfo: map[string][]byte{"xxx": []byte("22222")},
		// ContentType:  1,
		// CompressType: 1,
		CallType:   proto.CallType_CallService,
		SrcEntity:  &proto.Entity{ID: "222"},
		DstEntity:  &proto.Entity{ID: "111"},
		CheckFlags: codec.FlagDataChecksumIEEE,
	}
	hdr.Method = "new_client"
	pkt.SetHeader(hdr)
	pkt.WriteBody(nil)

	wg.Add(1)
	client.WritePacket(pkt)

	hdr.Method = "helloworld.Test/SayHello"
	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("raw data len:", len(pkt.Data()), len(hdrBytes))
	pkt.WriteBody(hdrBytes)
	client.WritePacket(pkt)

	wg.Wait()
}
