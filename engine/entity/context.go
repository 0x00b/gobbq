package entity

import "context"

type Context interface {
	context.Context

	GetEntityID() EntityID

	Parent() IEntity
	Childrens() []IEntity
}

type entityContext struct {
	context.Context

	entityID EntityID
}

// func (ec *entityContext) GetEntityID() EntityID {
// 	return ec.entityID
// }

// func WithEntityID(ctx context.Context, id EntityID) Context {
// 	return &entityContext{
// 		Context:  ctx,
// 		entityID: id,
// 	}
// }
