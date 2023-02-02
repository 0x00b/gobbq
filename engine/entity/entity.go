package entity

// just for inner
type EntityID string

// just for inner
type TypeName string

// IEntity declares functions that is defined in Entity
// These functions are mostly component functions
type IEntity interface {
	IService

	// entity 是有ID的service
	EntityID() EntityID
}

var _ IEntity = &Entity{}

type Entity struct {
	Service

	// parent maybe service or entity
	parent    IService
	childrens []IEntity
}

// entity 是有ID的service
func (e *Entity) EntityID() EntityID {
	return e.context.entityID
}

func (e *Entity) SetEntityID(id EntityID) {
	e.context.entityID = id
	return
}

// for inner

func (e *Entity) onDestroy() {
	e.OnDestroy()
	e.Service.onDestroy()
}
