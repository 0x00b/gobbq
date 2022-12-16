package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type ClientTransport struct {
	*ServerTransport
}

func NewClientTransport(ctx context.Context, conn net.Conn) *ClientTransport {
	ct := &ClientTransport{}
	ct.ServerTransport = NewServerTransportWithPacketHandler(ctx, conn, ct)
	return ct
}

func (ct *ClientTransport) HandlePacket(c context.Context, pkt *codec.Packet) error {

	fmt.Println("recv", string(pkt.PacketBody()))

	newpkt := codec.NewPacket()
	newpkt.WriteBytes([]byte("test"))

	err := ct.WritePacket(newpkt)
	if err != nil {
		fmt.Println("WritePacket", err)
	}

	return nil
}

func (ct *ClientTransport) Invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error {

	// req

	// register rsp use req.GetPacketID()

	return nil
}

func (ct *ClientTransport) invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error {

	var opts []CallOption //default opts

	ci := &CallInfo{
		Method: method,
	}

	pkt, err := ct.newPacket(ctx, ci, req, opts...)
	if err != nil {
		return err
	}
	if err := ct.WritePacket(pkt); err != nil {
		return err
	}
	return nil
}

// newPacket reffer parsePacket
func (ct *ClientTransport) newPacket(ctx context.Context, ci *CallInfo, req interface{}, opts ...CallOption) (*codec.Packet, error) {

	return nil, nil
}
