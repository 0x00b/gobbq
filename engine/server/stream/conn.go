package stream

import (
	"net"

	"github.com/0x00b/gobbq/engine/codec"
)

type conn struct {
	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter
	codec            codec.Codec
	compressor       codec.Compressor
}
