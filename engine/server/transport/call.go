package transport

import (
	"github.com/0x00b/gobbq/engine/codec"
)

type Header struct {
	messageID uint32
	flags     codec.Flags

	keys map[string][]string
}

// CallerHdr carries the information of a particular RPC.
type CallInfo struct {
	// Host specifies the peer's host.
	Host string

	// Method specifies the operation to perform.
	Method string

	Header Header
}

// CallerOption configures aCaller before it starts or extracts information from
// aCaller after it completes.
type CallOption interface {
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*CallInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*CallInfo)
}

// type CallerInterface interface {

// 	// NewMessage
// 	NewMessage(ctx context.Context, ci *CallInfo, req interface{}, opts ...CallOption) (*codec.Message, error)
// }

// type CalleeInterface interface {

// 	// ParseMessage
// 	ParseMessage(ctx context.Context, pkt *codec.Message) (ci *CallInfo, reqBodyBuff []byte, err error)
// }

// func Invoke(ctx context.Context, method string, req interface{}, callback ...func(ctx context.Context, reply interface{}) error) error
