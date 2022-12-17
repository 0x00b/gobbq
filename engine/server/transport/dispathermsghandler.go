package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type DispatherMessageHandler struct {
}

func NewDispatherTransport(ctx context.Context, conn net.Conn) *DispatherMessageHandler {
	st := &DispatherMessageHandler{}
	// st.ServerTransport = NewServerTransportWithMessageHandler(ctx, conn, st)
	return st
}

func (st *DispatherMessageHandler) HandleMessage(c context.Context, pkt *codec.Message) error {

	fmt.Println("recv", string(pkt.MessageBody()))
	// send to game
	// or send to gate

	return nil
}
