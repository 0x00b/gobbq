package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
)

func main() {

	fmt.Println(conf.C)

	var te TestServerInterface = &TestEntity{}

	RegisterTestService(te)

	RegisterTestEntity(te)

	// svr := gobbq.NewSever(nets.WithPacketHandler(game.NewGamePacketHandler()))
	// go svr.ListenAndServe(nets.TCP, ":1234")
	// go svr.ListenAndServe(nets.KCP, ":1235")
	// err := svr.ListenAndServe(nets.WebSocket, ":8080")
	game.ConnectProxy()

	testServer()

	bufio.NewReader(os.Stdin).ReadString('\n')
	// fmt.Println(err)
}

func testServer() {
	clientEty := NewTestEntityWithID("111")
	_ = clientEty

	rsp, err := clientEty.SayHello(context.Background(), &bbq.Header{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

type TestEntity struct {
	entity.Entity
}
type TestServerInterface interface {
	entity.IEntity
	TestInterface
}

type TestInterface interface {
	SayHello(c context.Context, req *bbq.Header) (*bbq.Header, error)
}

func (*TestEntity) SayHello(c context.Context, req *bbq.Header) (*bbq.Header, error) {

	return &bbq.Header{
		Method: "hello",
	}, nil
}

func RegisterTestEntity(svc TestServerInterface) {
	entity.Manager.RegisterEntity(&Test_ServiceDesc, svc)
}

func RegisterTestService(svc TestServerInterface) {
	entity.Manager.RegisterService(&Test_ServiceDesc, svc)
}

func _SayHello_Handler(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor entity.ServerInterceptor) {
	in := new(bbq.Header)
	if err := dec(in); err != nil {
		// return nil, err
		return
	}
	if interceptor == nil {
		svc.(*TestEntity).SayHello(ctx, in)
		return
	}
	return
}

var Test_ServiceDesc = entity.ServiceDesc{
	TypeName:    "helloworld.Test",
	HandlerType: (*TestServerInterface)(nil),
	Methods: map[string]entity.MethodDesc{
		"SayHello": {
			MethodName: "SayHello",
			// Handler:    _SayHello_Handler,
		},
	},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}

func NewTestEntity() *testClienEntity {
	return NewTestEntityWithID(entity.EntityID(snowflake.GenUUID()))
}

func NewTestEntityWithID(id entity.EntityID) *testClienEntity {

	ety := entity.NewEntity(id, Test_ServiceDesc.TypeName)
	t := &testClienEntity{entity: ety}

	return t
}

// client

func testClient() {
	cli := NewTestClient("id")
	rsp, err := cli.SayHello(context.Background(), &bbq.Header{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func NewTestClient(id entity.EntityID) *testClienEntity {
	c := testClienEntity{
		entity: &bbq.EntityID{
			ID:       string(id),
			TypeName: string(Test_ServiceDesc.TypeName),
		},
	}
	return &c
}

type testClienEntity struct {
	entity *bbq.EntityID
}

// 返回内容不可改
func (t *testClienEntity) GetEntity() *bbq.EntityID {
	e := &bbq.EntityID{}
	e.ID = t.entity.ID
	e.TypeName = t.entity.TypeName
	return e
}

func (t *testClienEntity) SayHello(c context.Context, req *bbq.Header) (*bbq.Header, error) {

	pkt := codec.NewPacket()

	hdr := &bbq.Header{
		Version:      1,
		RequestId:    "1",
		Timeout:      1,
		RequestType:  0,
		ServiceType:  0,
		SrcEntity:    t.entity,
		DstEntity:    t.entity,
		Method:       "helloworld.Test/SayHello",
		ContentType:  0,
		CompressType: 0,
		CheckFlags:   codec.FlagDataChecksumIEEE,
		TransInfo:    map[string][]byte{"xxx": []byte("22222")},
		ErrCode:      0,
		ErrMsg:       "",
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(bbq.ContentType_Proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return hdr, nil
	}

	pkt.WriteBody(hdrBytes)

	ex.SendProxy(pkt)

	return hdr, nil
}
