package game

import (
	"fmt"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
)

type RegisterProxy struct {
}

func (*RegisterProxy) RegisterEntityToProxy(eid entity.EntityID) error {
	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterEntity(nil, &proxypb.RegisterEntityRequest{EntityID: string(eid)},
		func(c *entity.Context, rsp *proxypb.RegisterEntityResponse) {
			fmt.Println("register proxy entity resp")
		},
	)

	return nil
}

func (*RegisterProxy) RegisterServiceToProxy(svcName entity.TypeName) error {

	client := proxypb.NewProxyServiceClient(ex.ProxyClient)

	client.RegisterService(nil, &proxypb.RegisterServiceRequest{ServiceName: string(svcName)},
		func(c *entity.Context, rsp *proxypb.RegisterServiceResponse) {

			fmt.Println("register proxy service resp")
		},
	)

	return nil
}
