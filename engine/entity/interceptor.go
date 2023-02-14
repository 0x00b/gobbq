package entity

import "github.com/0x00b/gobbq/engine/codec"

type ServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server any
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

// 请求回调
type Callback func(pkt *codec.Packet)

type Handler func(ctx Context, req any) (any, error)

type ServerInterceptor func(ctx Context, req any, info *ServerInfo, next Handler) (any, error)

// chainServerInterceptors chains all  server interceptors into one.
func chainServerInterceptors(interceptors []ServerInterceptor) ServerInterceptor {
	// Prepend opts.Int to the chaining interceptors if it exists, since Int will
	// be executed before any other chained interceptors.

	var chainedInt ServerInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = chainInterceptors(interceptors)
	}

	return chainedInt
}

func chainInterceptors(interceptors []ServerInterceptor) ServerInterceptor {
	return func(ctx Context, req any, info *ServerInfo, handler Handler) (any, error) {
		return interceptors[0](ctx, req, info, getChainHandler(interceptors, 0, info, handler))
	}
}

func getChainHandler(interceptors []ServerInterceptor, curr int, info *ServerInfo, finalHandler Handler) Handler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx Context, req any) (any, error) {
		return interceptors[curr+1](ctx, req, info, getChainHandler(interceptors, curr+1, info, finalHandler))
	}
}
