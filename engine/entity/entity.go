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
	SetEntityID(id EntityID)

	// 有状态

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in
	// Freeze && Restore
	OnFreeze()   // Called when entity is freezing
	OnRestored() // Called when entity is restored

}

var _ IEntity = &Entity{}

type Entity struct {
	Service

	entityID EntityID
}

// entity 是有ID的service
func (e *Entity) EntityID() EntityID {
	return e.entityID
}

func (e *Entity) SetEntityID(id EntityID) {
	e.entityID = id
	return
}

// Migration
func (e *Entity) OnMigrateOut() {} // Called just before entity is migrating out
func (e *Entity) OnMigrateIn()  {} // Called just after entity is migrating in
// Freeze && Restore
func (e *Entity) OnFreeze()   {} // Called when entity is freezing
func (e *Entity) OnRestored() {} // Called when entity is restored

type EntityClient struct {
	// todo
	Entity
}
