package ex

import (
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/nets"
)

var ProxyClient *nets.Client

func ConnProxy(cfg conf.NetConf, opts ...nets.Option) {

	cli, err := nets.Connect(nets.NetWorkName(cfg.Net), cfg.IP, cfg.Port, opts...)

	if err != nil {
		panic(err)
	}

	ProxyClient = cli
}

func SendProxy(pkt *nets.Packet) error {

	return ProxyClient.SendPacket(pkt)

}
