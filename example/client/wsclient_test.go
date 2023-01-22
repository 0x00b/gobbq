package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto"
	"golang.org/x/net/websocket"
)

func TestWSClient(m *testing.T) {

	origin := "http://localhost:8080/"
	url := "ws://localhost:8080/ws"
	wsc, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	ws := codec.NewPacketReadWriter(context.Background(), wsc)

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "helloworld.Test/SayHello",
		TransInfo: map[string][]byte{"xxx": []byte("22222")},
		// ContentType:  1,
		// CompressType: 1,
		DstEntity:  &proto.Entity{ID: "Y80_q1ZNLX9eAAAB"},
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}

	pkt.WriteBody(hdrBytes)

	fmt.Println("raw data len:", len(pkt.Data()), len(hdrBytes))

	ws.WritePacket(pkt)

	if pkt, err = ws.ReadPacket(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", string(pkt.PacketBody()))
}
