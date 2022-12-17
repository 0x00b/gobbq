package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type ClientTransport struct {
	*conn
}

func NewClientTransport(ctx context.Context, rawConn net.Conn) *ClientTransport {
	ct := &ClientTransport{
		conn: &conn{
			rwc:               rawConn,
			ctx:               context.Background(),
			messageReadWriter: codec.NewMessageReadWriter(context.Background(), rawConn),
			MessageHandler:    NewServerMessageHandler(context.Background(), rawConn),
		},
	}
	// ct.ServerTransport = NewServerTransportWithMessageHandler(ctx, conn, ct)
	return ct
}

func (ct *ClientTransport) HandleMessage(c context.Context, pkt *codec.Message) error {

	fmt.Println("recv", string(pkt.MessageBody()))

	newpkt := codec.NewMessage()
	newpkt.WriteBytes([]byte("test"))

	err := ct.WriteMessage(newpkt)
	if err != nil {
		fmt.Println("WriteMessage", err)
	}

	return nil
}

func (ct *ClientTransport) Invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error {

	// req

	// register rsp use req.GetMessageID()

	return nil
}

func (ct *ClientTransport) invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error {

	var opts []CallOption //default opts

	ci := &CallInfo{
		Method: method,
	}

	pkt, err := ct.newMessage(ctx, ci, req, opts...)
	if err != nil {
		return err
	}
	if err := ct.WriteMessage(pkt); err != nil {
		return err
	}
	return nil
}

// newMessage reffer parseMessage
func (ct *ClientTransport) newMessage(ctx context.Context, ci *CallInfo, req interface{}, opts ...CallOption) (*codec.Message, error) {

	return nil, nil
}
