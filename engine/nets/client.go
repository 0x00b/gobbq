package nets

import (
	"context"
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type Client struct {
	*conn
}

func NewClient(ctx context.Context, rawConn net.Conn, ops ...Option) *Client {

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
	return ct.conn.WritePacket(pkt)
}
