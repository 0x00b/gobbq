// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package codec

import (
	"fmt"

	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/golang/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() {
	RegisterCodec(DefaultCodec)
}

var DefaultCodec = protoCodec{}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type protoCodec struct{}

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, packet is %T, want proto.Packet", v)
	}
	return proto.Marshal(vv)
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, packet is %T, want proto.Packet", v)
	}
	return proto.Unmarshal(data, vv)
}

func (protoCodec) Type() bbq.ContentType {
	return bbq.ContentType_Proto
}
