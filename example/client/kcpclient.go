package main

import (
	"context"
	"fmt"
	"time"

	"github.com/0x00b/gobbq/bbqpb"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server/transport"
	"github.com/xtaci/kcp-go"
)

func main() {

	wsc, err := kcp.DialWithOptions("127.0.0.1:1235", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("runing")

	ctx := context.Background()

	ct := transport.NewClientTransport(ctx, wsc)

	pkt := codec.NewPacket()

	hdr := &bbqpb.RequestHeader{
		Version:   1,
		RequestId: 1,
		Timeout:   1,
		Method:    "helloworld.Test/SayHello",
		TransInfo: map[string][]byte{"xxx": []byte("22222")},
		// ContentType:  1,
		// CompressType: 1,
	}

	hdrBytes, err := codec.GetCodec(bbqpb.ContentType_proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("len:", pkt.GetMsgHeaderLen(), uint32(len("dsfsdfs")), pkt.GetPacketBodyLen())
	fmt.Println("raw data:", []byte(hdr.String()), []byte("dsfsdfs"))

	pkt.WriteMsgHeader(hdrBytes)

	// body
	pkt.WriteBytes(hdrBytes)

	fmt.Println("data:", len(pkt.PacketBody()), pkt.PacketBody())

	ct.WritePacket(pkt)

	ct.Serve()
	time.Sleep(1 * time.Minute)
}
