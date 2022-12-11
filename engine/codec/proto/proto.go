// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"fmt"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/golang/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() {
	codec.RegisterCodec(protoCodec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type protoCodec struct{}

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return proto.Marshal(vv)
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}

func (protoCodec) Name() string {
	return Name
}
