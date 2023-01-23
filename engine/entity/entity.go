package entity

import (
	"github.com/0x00b/gobbq/proto"
)

// just for inner
type EntityID string

// just for inner
type ServiceType string

// IEntity declares functions that is defined in Entity
// These functions are mostly component functions
type IEntity interface {
	IService

	// entity 是有ID的service
	Entity() *proto.Entity
	// entity 是有ID的service
	SetEntityID(id EntityID)

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in
	// Freeze && Restore
	OnFreeze()   // Called when entity is freezing
	OnRestored() // Called when entity is restored

}

var _ IEntity = &NopEntity{}

type NopEntity struct {
	NopService

	ety *proto.Entity
}

// entity 是有ID的service
func (n *NopEntity) Entity() *proto.Entity {
	return n.ety
}

func (n *NopEntity) SetEntityID(id EntityID) {
	if n.ety == nil {
		n.ety = &proto.Entity{
			ID:   string(id),
			Type: "Entity",
		}

	}
	return
}

// Migration
func (n *NopEntity) OnMigrateOut() {} // Called just before entity is migrating out
func (n *NopEntity) OnMigrateIn()  {} // Called just after entity is migrating in
// Freeze && Restore
func (n *NopEntity) OnFreeze()   {} // Called when entity is freezing
func (n *NopEntity) OnRestored() {} // Called when entity is restored
