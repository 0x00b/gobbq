package main

import (
	"context"
	"fmt"

	"github.com/0x00b/gobbq"
	"github.com/0x00b/gobbq/bbqpb"
	"github.com/0x00b/gobbq/engine/server"
)

func main() {
	svr := gobbq.NewSever()

	RegisterTestEntity(svr, &TestEntity{})

	go svr.ListenAndServe(server.TCP, ":1234")
	go svr.ListenAndServe(server.KCP, ":1235")
	err := svr.ListenAndServe(server.WebSocket, ":80")

	fmt.Println(err)
}

type TestEntity struct {
}

type TestEntityInterface interface {
	SayHello(c context.Context, req *bbqpb.RequestHeader) (*bbqpb.ResponseHeader, error)
}

func (*TestEntity) SayHello(c context.Context, req *bbqpb.RequestHeader) (*bbqpb.ResponseHeader, error) {
	return &bbqpb.ResponseHeader{
		Method: "hello",
	}, nil
}

func RegisterTestEntity(s *server.Server, srv TestEntityInterface) {
	s.RegisterEntity(&Test_ServiceDesc, srv)
}

func _SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.UnaryServerInterceptor) (interface{}, error) {
	in := new(bbqpb.RequestHeader)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*TestEntity).SayHello(ctx, in)
	}
	return nil, nil
}

var Test_ServiceDesc = server.EntityDesc{
	TypeName:    "helloworld.Test",
	HandlerType: (*TestEntityInterface)(nil),
	Methods: map[string]server.MethodDesc{
		"SayHello": {
			MethodName: "SayHello",
			Handler:    _SayHello_Handler,
		},
	},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}
