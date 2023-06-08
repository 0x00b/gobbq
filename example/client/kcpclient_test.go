package client_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/xtaci/kcp-go"
)

func TestMain(t *testing.T) {

	wsc, err := kcp.DialWithOptions("127.0.0.1:8899", nil, 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("runing")

	pkt, _ := codec.NewPacket()

	hdr := &bbq.Header{
		Version:   1,
		RequestId: "1",
		Timeout:   1,
		Method:    "helloworld.Test/SayHello",
		TransInfo: map[string][]byte{"xxx": []byte("22222")},
		// ContentType:  1,
		// CompressType: 1,
	}

	pkt.Header = (hdr)

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("raw data:", []byte(hdr.String()), []byte("dsfsdfs"))

	// body
	pkt.WriteBody(hdrBytes)

	fmt.Println("data:", len(pkt.PacketBody()), pkt.PacketBody())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		i := 0
		for {
			var b [1024]byte
			n, _ := wsc.Read(b[:])
			fmt.Println("recv", i, string(b[:n]))
			i++
		}
	}()

	for i := 0; i < 1000; i++ {
		fmt.Println("send", i, string(pkt.Data()))
		wsc.Write(pkt.Data())
	}

	wg.Wait()

}
