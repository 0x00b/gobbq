package main

func RegisterEntityToProxy(eid string) error {

	// ps := NewProxyService(nil)

	// ps.RegisterEntity(nil, &RegisterEntityRequest{
	// 	EntityID: eid,
	// }, func(c *entity.Context, rsp *RegisterEntityResponse) {
	// 	fmt.Println("recv:", string(c.Packet().PacketBody()))
	// 	fmt.Println(rsp)
	// })

	// pkt, release := codec.NewPacket()
	// defer release()

	// pkt.Header.Version = 1
	// pkt.Header.RequestId = "1"
	// pkt.Header.Timeout = 1
	// pkt.Header.RequestType = 0
	// pkt.Header.ServiceType = 0
	// pkt.Header.SrcEntity = &bbq.EntityID{ID: eid}
	// pkt.Header.DstEntity = &bbq.EntityID{}
	// pkt.Header.Method = "register_proxy_entity"
	// pkt.Header.ContentType = 0
	// pkt.Header.CompressType = 0
	// pkt.Header.CheckFlags = codec.FlagDataChecksumIEEE
	// pkt.Header.TransInfo = map[string][]byte{"xxx": []byte("22222")}
	// pkt.Header.ErrCode = 0
	// pkt.Header.ErrMsg = ""

	// pkt.WriteBody(nil)

	// fmt.Println("register", string(pkt.Header.GetSrcEntity().ID))

	// return proxyMap[1].SendPackt(pkt)

	return nil
}
