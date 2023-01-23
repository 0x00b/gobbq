package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto"
)

func main() {
	svr := gobbq.NewSever(nets.WithPacketHandler(game.NewGamePacketHandler()))

	var te TestEntityInterface = &TestEntity{}

	RegisterTestService(svr, te)

	RegisterTestEntity(svr, te)

	entity.NewEntity(Test_ServiceDesc.TypeName)

	go svr.ListenAndServe(nets.TCP, ":1234")
	go svr.ListenAndServe(nets.KCP, ":1235")
	err := svr.ListenAndServe(nets.WebSocket, ":8080")

	fmt.Println(err)
}

type TestEntity struct {
	entity.NopEntity
}

type TestEntityInterface interface {
	entity.IEntity
	SayHello(c context.Context, req *proto.Header) (*proto.Header, error)
}

func (*TestEntity) SayHello(c context.Context, req *proto.Header) (*proto.Header, error) {
	return &proto.Header{
		Method: "hello",
	}, nil
}

func RegisterTestEntity(s *nets.Server, svc TestEntityInterface) {
	entity.Manager.RegisterEntity(&Test_ServiceDesc, svc)
}

func RegisterTestService(s *nets.Server, svc TestEntityInterface) {
	entity.Manager.RegisterService(&Test_ServiceDesc, svc)
}

func _SayHello_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.UnaryServerInterceptor) (interface{}, error) {
	in := new(proto.Header)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return svc.(*TestEntity).SayHello(ctx, in)
	}
	return nil, nil
}

var Test_ServiceDesc = entity.ServiceDesc{
	TypeName:    "helloworld.Test",
	HandlerType: (*TestEntityInterface)(nil),
	Methods: map[string]entity.MethodDesc{
		"SayHello": {
			MethodName: "SayHello",
			Handler:    _SayHello_Handler,
		},
	},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}
