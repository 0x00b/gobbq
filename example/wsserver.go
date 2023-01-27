package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0x00b/gobbq/components/game"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto"
	"golang.org/x/net/websocket"
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
	// rsp, err := clientEty.SayHello(context.Background(), &proto.Header{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(rsp)
}

type TestEntity struct {
	entity.Entity
}
type TestServerInterface interface {
	entity.IEntity
	TestInterface
}

type TestInterface interface {
	SayHello(c context.Context, req *proto.Header) (*proto.Header, error)
}

func (*TestEntity) SayHello(c context.Context, req *proto.Header) (*proto.Header, error) {
	return &proto.Header{
		Method: "hello",
	}, nil
}

func RegisterTestEntity(svc TestServerInterface) {
	entity.Manager.RegisterEntity(&Test_ServiceDesc, svc)
}

func RegisterTestService(svc TestServerInterface) {
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
	HandlerType: (*TestServerInterface)(nil),
	Methods: map[string]entity.MethodDesc{
		"SayHello": {
			MethodName: "SayHello",
			Handler:    _SayHello_Handler,
		},
	},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}

// func NewTestEntity() *testClienEntity {
// 	return NewTestEntityWithID(entity.EntityID(snowflake.GenUUID()))
// }

func NewTestEntityWithID(id entity.EntityID) *testClienEntity {

	ety := entity.NewEntity(id, Test_ServiceDesc.TypeName)
	t := &testClienEntity{entity: ety}

	return t
}

// client

func testClient() {
	cli := NewTestClient("id")
	rsp, err := cli.SayHello(context.Background(), &proto.Header{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func NewTestClient(id entity.EntityID) TestInterface {
	c := testClienEntity{
		entity: &proto.Entity{
			ID:       string(id),
			TypeName: string(Test_ServiceDesc.TypeName),
		},
	}
	return &c
}

type testClienEntity struct {
	entity *proto.Entity
}

// 返回内容不可改
func (t *testClienEntity) GetEntity() *proto.Entity {
	e := &proto.Entity{}
	e.ID = t.entity.ID
	e.TypeName = t.entity.TypeName
	return e
}

func (t *testClienEntity) SayHello(c context.Context, req *proto.Header) (*proto.Header, error) {

	origin := "http://localhost:8080/"
	url := "ws://localhost:8080/ws"
	wsc, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	ws := codec.NewPacketReadWriter(context.Background(), wsc)

	pkt := codec.NewPacket()

	hdr := &proto.Header{
		Version:    1,
		RequestId:  "1",
		Timeout:    1,
		Method:     "helloworld.Test/SayHello",
		TransInfo:  map[string][]byte{"xxx": []byte("22222")},
		CallType:   proto.CallType_CallService,
		SrcEntity:  t.entity,
		DstEntity:  t.entity,
		CheckFlags: codec.FlagDataChecksumIEEE,
	}

	pkt.SetHeader(hdr)

	hdrBytes, err := codec.GetCodec(proto.ContentType_Proto).Marshal(hdr)
	if err != nil {
		fmt.Println(err)
		return hdr, nil
	}

	pkt.WriteBody(hdrBytes)

	fmt.Println("raw data len:", len(pkt.Data()), len(hdrBytes))

	// todo proxy send
	ws.WritePacket(pkt)

	if pkt, err = ws.ReadPacket(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", string(pkt.PacketBody()))

	return hdr, nil
}
