package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type GateMessageHandler struct {
}

func NewGateMessageHandler(ctx context.Context, conn net.Conn) *GateMessageHandler {
	st := &GateMessageHandler{}
	// st.ServerTransport = NewServerTransportWithMessageHandler(ctx, conn, st)
	return st
}

func (st *GateMessageHandler) HandleMessage(c context.Context, pkt *codec.Message) error {

	fmt.Println("recv", string(pkt.MessageBody()))

	// send to dispather
	return nil
}
