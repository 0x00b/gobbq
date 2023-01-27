package ex

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto"
)

type ProxyMap map[uint32]*nets.Client

var proxyMap ProxyMap = make(ProxyMap)

func ConnProxy(ops ...nets.Option) {
	for i := 0; i < int(conf.C.Proxy.InstNum); i++ {

		// connect to proxy
		cfg := conf.C.Proxy.Inst[i]
		_ = cfg.Net

		prxy, err := nets.Connect(context.Background(), nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, ops...)

		if err != nil {
			panic(err)
		}

		proxyMap[cfg.ID] = prxy
	}
}

func SendProxy(pkt *codec.Packet) error {
	_ = pkt.GetHeader().GetDstEntity().ID
	// hash id , lb proxy

	return proxyMap[1].SendPackt(pkt)

}

func RegisterEntity(eid string) error {

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:    1,
		RequestId:  "1",
		Timeout:    1,
		Method:     "register_proxy_entity",
		TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallService,
		SrcEntity:  &proto.Entity{ID: eid},
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)
	pkt.WriteBody(nil)

	fmt.Println("register", string(hdr.GetSrcEntity().ID))

	return proxyMap[1].SendPackt(pkt)

}
