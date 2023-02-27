package main

func main() {

	// wsc, err := kcp.DialWithOptions("127.0.0.1:1235", nil, 10, 3)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("runing")

	// ctx := context.Background()

	// ct := transport.NewClientTransport(ctx, wsc)

	// pkt := codec.NewPacket()

	// hdr := &bbq.Header{
	// 	Version:   1,
	// 	RequestId: "1",
	// 	Timeout:   1,
	// 	Method:    "helloworld.Test/SayHello",
	// 	TransInfo: map[string][]byte{"xxx": []byte("22222")},
	// 	// ContentType:  1,
	// 	// CompressType: 1,
	// }

	// pkt.SetHeader(hdr)

	// hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(hdr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("raw data:", []byte(hdr.String()), []byte("dsfsdfs"))

	// // body
	// pkt.WriteBody(hdrBytes)

	// fmt.Println("data:", len(pkt.PacketBody()), pkt.PacketBody())

	// ct.SendPackt(pkt)

	// ct.Serve()
	// time.Sleep(1 * time.Minute)
}
