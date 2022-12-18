package server

import (
	"context"
)

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

// EntityDesc represents an RPC service's specification.
type EntityDesc struct {
	ServiceImpl interface{}
	TypeName    string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     map[string]MethodDesc
	Metadata    interface{}
}

// ClientID type
type ClientID string

// GameClient represents the game Client of Entity
//
// Each Entity can have at most one GameClient, and GameClient can be given to other Entitys
type GameClient struct {
	clientID ClientID
	gateID   uint16
	ownerID  string
}
