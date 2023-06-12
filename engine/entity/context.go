package entity

import (
	"context"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

type Context interface {

	/************************************/
	/***** GOLANG.ORG/X/NET/CONTEXT *****/
	/************************************/
	context.Context

	/************************************/
	/******** ENTITY MANAGEMENT********/
	/************************************/
	Copy() (Context, releaseCtx)

	Entity() IBaseEntity

	EntityID() EntityID

	Packet() *codec.Packet

	SrcEntity() EntityID

	/************************************/
	/******** METADATA MANAGEMENT********/
	/************************************/
	SetError(err error)

	// Set is used to store a new key/value pair exclusively for this context.
	// It also lazy initializes  c.Keys if it was not used previously.
	Set(key string, value any)

	// Get returns the value for the given key, ie: (value, true).
	// If the value does not exist it returns (nil, false)
	Get(key string) (value any, exists bool)

	// MustGet returns the value for the given key if it exists, otherwise it panics.
	MustGet(key string) any

	// GetString returns the value associated with the key as a string.
	GetString(key string) (s string)

	// GetBool returns the value associated with the key as a boolean.
	GetBool(key string) (b bool)

	// GetInt returns the value associated with the key as an integer.
	GetInt(key string) (i int)

	// GetInt64 returns the value associated with the key as an integer.
	GetInt64(key string) (i64 int64)

	// GetUint returns the value associated with the key as an unsigned integer.
	GetUint(key string) (ui uint)

	// GetUint64 returns the value associated with the key as an unsigned integer.
	GetUint64(key string) (ui64 uint64)

	// GetFloat64 returns the value associated with the key as a float64.
	GetFloat64(key string) (f64 float64)
	// GetTime returns the value associated with the key as time.
	GetTime(key string) (t time.Time)

	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) (d time.Duration)

	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) (ss []string)

	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) (sm map[string]any)

	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(key string) (sms map[string]string)

	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(key string) (smss map[string][]string)

	// inner
	setPacket(pkt *codec.Packet)
}

// ============ for bbq inner start=================
// 不想通过context暴露给其他开发者

const (
	// _bbq_ctx_pkt_key_            = "_bbq_ctx_pkt_key_"
	_bbq_ctx_remote_manager_key_ = "_bbq_ctx_remote_manager_key_"
	// _bbq_ctx_entity_manager_key_ = "_bbq_ctx_entity_manager_key_"
)

func RegisterCallback(c Context, requestID string, cb Callback) {
	e := c.Entity()
	if e != nil {
		e.registerCallback(requestID, cb)
	}
}

func PopCallback(c Context, requestID string) (Callback, bool) {
	e := c.Entity()
	if e != nil {
		return e.popCallback(requestID)
	}
	return nil, false
}

func GetProxy(c Context) Proxy {
	v, ok := c.Get(_bbq_ctx_remote_manager_key_)
	if ok && v != nil {
		return v.(Proxy)
	}

	etyMgr := GetEntityMgr(c)
	return etyMgr.Proxy
}

func SetRemoteEntityManager(c Context, rem Proxy) {
	if rem == nil {
		return
	}
	c.Set(_bbq_ctx_remote_manager_key_, rem)
}

func GetEntityMgr(c Context) *EntityManager {
	// v, ok := c.Get(_bbq_ctx_entity_manager_key_)
	// if ok && v != nil {
	// 	return v.(*EntityManager)
	// }
	return c.Entity().getEntityMgr()
}

// func SetEntityMgr(c Context, em *EntityManager) {
// 	if em == nil {
// 		return
// 	}
// 	c.Set(_bbq_ctx_entity_manager_key_, em)
// }

// ============ for bbq inner end=================

var _ Context = &baseContext{}

var contextPool *sync.Pool = &sync.Pool{
	New: func() any {
		c := &baseContext{}
		return c
	},
}

func allocContext(parent context.Context) (*baseContext, func()) {
	c := contextPool.Get().(*baseContext)

	c.reset()
	if parent == nil {
		parent = context.Background()
	}
	tctx, cancel := context.WithCancel(parent)
	c.ctx = tctx

	return c, cancel
}

func releaseContext(c Context) {
	if c != nil {
		bc, ok := c.(*baseContext)
		if ok {
			bc.reset()
			contextPool.Put(bc)
		}
	}
}

type releaseCtx func()

// NewPacket allocates a new packet
func NewContext(parent context.Context) (Context, releaseCtx) {
	c, cancel := allocContext(parent)

	return c, func() {
		cancel()
		releaseContext(c)
	}
}

type baseContext struct {
	ctx context.Context

	// 属于这个entity
	entity IBaseEntity

	err error

	// This mutex protects Keys map.
	mu sync.RWMutex
	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any

	pkt *codec.Packet
}

func (c *baseContext) Entity() IBaseEntity {
	return c.entity
}

func (c *baseContext) Packet() *codec.Packet {
	return c.pkt
}

func (c *baseContext) setPacket(pkt *codec.Packet) {
	c.pkt = pkt
}

func (c *baseContext) SrcEntity() EntityID {
	return EntityID(c.pkt.Header.GetSrcEntity())
}

func (c *baseContext) Copy() (Context, releaseCtx) {

	tc, cancel := allocContext(c.ctx)
	tc.reset()

	tc.entity = c.entity
	tc.err = c.err
	tc.mu = sync.RWMutex{}
	tc.Keys = map[string]any{}

	for k, v := range c.Keys {
		tc.Set(k, v)
	}

	return tc, func() {
		cancel()
		releaseContext(tc)
	}
}

func (c *baseContext) reset() {
	c.entity = nil
	c.ctx = nil
	c.err = nil
}

func (c *baseContext) EntityID() EntityID {
	return c.entity.EntityID()
}

func (c *baseContext) SetError(err error) {
	if err == nil {
		panic("err is nil")
	}
	c.err = err
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *baseContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *baseContext) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *baseContext) MustGet(key string) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *baseContext) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *baseContext) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *baseContext) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *baseContext) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *baseContext) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *baseContext) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *baseContext) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *baseContext) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *baseContext) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *baseContext) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *baseContext) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *baseContext) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *baseContext) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline returns that there is no deadline (ok==false) when c.ctx has no Context.
func (c *baseContext) Deadline() (deadline time.Time, ok bool) {
	if c.ctx == nil {
		return
	}
	return c.ctx.Deadline()
}

// Done returns nil (chan which will wait forever) when c.ctx has no Context.
func (c *baseContext) Done() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Done()
}

// Err returns nil when c.ctx has no Context.
func (c *baseContext) Err() error {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *baseContext) Value(key any) any {
	if key == 0 {
		return c.ctx
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Value(key)
}
