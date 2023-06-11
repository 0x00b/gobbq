package ex

import (
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
)

var ProxyClient *nets.Client

func ConnProxy(ops ...nets.Option) {

	cfg := conf.C.Proxy.Inst[0]

	cli, err := nets.Connect(nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, ops...)

	if err != nil {
		panic(err)
	}

	ProxyClient = cli
}

func SendProxy(pkt *codec.Packet) error {

	return ProxyClient.SendPacket(pkt)

}
