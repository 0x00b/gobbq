package entity

import (
	"context"
)

type ServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server interface{}
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

type RetFunc func(interface{}, error)

type Handler func(ctx context.Context, req interface{}, ret RetFunc)

type ServerInterceptor func(ctx context.Context, req interface{}, info *ServerInfo, ret RetFunc, next Handler)

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
	return func(ctx context.Context, req interface{}, info *ServerInfo, ret RetFunc, handler Handler) {
		interceptors[0](ctx, req, info, ret, getChainHandler(interceptors, 0, info, ret, handler))
		return
	}
}

func getChainHandler(interceptors []ServerInterceptor, curr int, info *ServerInfo, ret RetFunc, finalHandler Handler) Handler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, req interface{}, ret RetFunc) {
		interceptors[curr+1](ctx, req, info, ret, getChainHandler(interceptors, curr+1, info, ret, finalHandler))
		return
	}
}
