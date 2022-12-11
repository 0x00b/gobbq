package entity

type EntityID string

type EntityType string

// IEntity declares functions that is defined in Entity
// These functions are mostly component functions
type IEntity interface {
	EntityID() EntityID

	// Entity Lifetime
	OnInit()       // Called when initializing entity struct, override to initialize entity custom fields
	OnAttrsReady() // Called when entity attributes are ready.
	OnCreated()    // Called when entity is just created
	OnDestroy()    // Called when entity is destroying (just before destroy)
	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in
	// Freeze && Restore
	OnFreeze()   // Called when entity is freezing
	OnRestored() // Called when entity is restored

	// Type returns the name of the Entity implementation.
	// the result cannot change between calls.
	Type() EntityType
}

// Entity is the basic execution unit in GoWorld server. Entities can be used to
// represent players, NPCs, monsters. Entities can migrate among spaces.
type Entity struct {
	I  IEntity
	ID EntityID

	typeName EntityType

	destroyed bool

	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	entityInfo *EntityDesc

	// syncingFromClient bool
	// rawTimers            map[*timer.Timer]struct{}
	// timers               map[EntityTimerID]*entityTimerInfo
	// lastTimerId          EntityTimerID
	// Attrs                *MapAttr
	// syncInfoFlag         syncInfoFlag
}
