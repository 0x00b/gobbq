package ex

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
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
	_ = pkt.Header.GetDstEntity().ID
	// hash id , lb proxy

	return proxyMap[1].SendPackt(pkt)

}

func RegisterEntity(eid string) error {

	pkt, release := codec.NewPacket()
	defer release()

	pkt.Header.Version = 1
	pkt.Header.RequestId = "1"
	pkt.Header.Timeout = 1
	pkt.Header.RequestType = 0
	pkt.Header.ServiceType = 0
	pkt.Header.SrcEntity = &bbq.EntityID{ID: eid}
	pkt.Header.DstEntity = &bbq.EntityID{}
	pkt.Header.Method = "register_proxy_entity"
	pkt.Header.ContentType = 0
	pkt.Header.CompressType = 0
	pkt.Header.CheckFlags = codec.FlagDataChecksumIEEE
	pkt.Header.TransInfo = map[string][]byte{"xxx": []byte("22222")}
	pkt.Header.ErrCode = 0
	pkt.Header.ErrMsg = ""

	pkt.WriteBody(nil)

	fmt.Println("register", string(pkt.Header.GetSrcEntity().ID))

	return proxyMap[1].SendPackt(pkt)

}
