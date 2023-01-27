package nets

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/xtaci/kcp-go"
	"golang.org/x/net/websocket"
)

type Client struct {
	*conn
}

func Connect(ctx context.Context, netname NetWorkName, ip, port string, ops ...Option) (*Client, error) {

	var conn net.Conn
	var err error
	switch netname {
	case KCP:
		conn, err = connectKcp(ctx, ip, port, ops...)
	case TCP, TCP6:
		conn, err = connectTcp(ctx, ip, port, ops...)
	case WebSocket:
		conn, err = connectWebsocket(ctx, ip, port, ops...)
	default:
		return nil, errors.New("unknown network")
	}
	if err != nil {
		return nil, err
	}

	return newClient(ctx, conn, ops...), nil
}

func connectTcp(ctx context.Context, ip, port string, ops ...Option) (net.Conn, error) {
	rc, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	return rc, nil
}

func connectKcp(ctx context.Context, ip, port string, ops ...Option) (net.Conn, error) {
	rc, err := kcp.DialWithOptions(fmt.Sprintf("%s:%s", ip, port), nil, 10, 3)
	if err != nil {
		panic(err)
	}
	return rc, nil
}

func connectWebsocket(ctx context.Context, ip, port string, ops ...Option) (net.Conn, error) {
	origin := fmt.Sprintf("http://%s:%s/", ip, port)
	url := fmt.Sprintf("ws://%s:%s/ws", ip, port)
	rc, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	return rc, nil
}

func newClient(ctx context.Context, rawConn net.Conn, ops ...Option) *Client {

	opts := &Options{}

	for _, op := range ops {
		op.apply(opts)
	}

	ct := &Client{
		conn: &conn{
			rwc:              rawConn,
			ctx:              ctx,
			packetReadWriter: codec.NewPacketReadWriter(ctx, rawConn),
			PacketHandler:    opts.PacketHandler,
			opts:             opts,
		},
	}

	go ct.conn.Serve()

	return ct
}

func (ct *Client) SendPackt(pkt *codec.Packet) error {
	// todo chan
	return ct.conn.WritePacket(pkt)
}
