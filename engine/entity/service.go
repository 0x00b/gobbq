package entity

import (
	"context"
)

type IService interface {
	// Entity Lifetime
	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnDestroy() // Called when entity is destroying (just before destroy)

	TypeName() TypeName
}

type methodHandler func(svc interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct {
	ServiceImpl interface{}

	TypeName TypeName
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     map[string]MethodDesc
	Metadata    interface{}
}

var _ IService = &Service{}

type Service struct {
	typeName TypeName
}

func (*Service) OnInit() {
	// Called when initializing entity struct, override to initialize entity custom fields
}

func (*Service) OnDestroy() {
	// Called when entity is destroying (just before destroy)
}

func (s *Service) TypeName() TypeName {
	return s.typeName
}
