package entity

import "context"

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

// EntityDesc represents an RPC service's specification.
type EntityDesc struct {
	TypeName EntityType
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     []MethodDesc
	Metadata    interface{}
}

// EntityInfo wraps information about a service. It is very similar to
// EntityDesc and is constructed from it for internal purposes.
type EntityInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*MethodDesc
	mdata       interface{}
}
