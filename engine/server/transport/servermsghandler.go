package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type ServerMessageHandler struct {
}

func NewServerMessageHandler(ctx context.Context, conn net.Conn) *ServerMessageHandler {
	st := &ServerMessageHandler{}
	// st.ServerTransport = NewServerTransportWithMessageHandler(ctx, conn, st)
	return st
}

func (st *ServerMessageHandler) HandleMessage(c context.Context, pkt *codec.Message) error {

	fmt.Println("recv", string(pkt.MessageBody()))

	newpkt := codec.NewMessage()
	newpkt.WriteBytes([]byte("test"))

	err := pkt.Src.WriteMessage(newpkt)
	if err != nil {
		fmt.Println("WriteMessage", err)
	}

	return nil
}
