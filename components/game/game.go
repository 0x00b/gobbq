package game

import (
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
)

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid entity.EntityID) error {
	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterEntity(nil, &proxypb.RegisterEntityRequest{EntityID: string(eid)},
		func(c *entity.Context, rsp *proxypb.RegisterEntityResponse) {
			xlog.Println("register proxy entity resp")
		},
	)

	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterService(nil, &proxypb.RegisterServiceRequest{ServiceName: string(svcName)},
		func(c *entity.Context, rsp *proxypb.RegisterServiceResponse) {

			xlog.Println("register proxy service resp")
		},
	)

	return nil
}
