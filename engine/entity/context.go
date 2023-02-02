package entity

import (
	"context"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
)

var _ context.Context = &Context{}

var contextPool *sync.Pool = &sync.Pool{
	New: func() any {
		c := &Context{}
		return c
	},
}

func allocContext() *Context {
	c := contextPool.Get().(*Context)

	c.reset()

	return c
}

func releaseContext(c *Context) {
	if c != nil {
		c.reset()
		contextPool.Put(c)
	}
}

type releaseCtx func()

// NewPacket allocates a new packet
func NewContext() (*Context, releaseCtx) {
	c := allocContext()
	return c, func() {
		releaseContext(c)
	}
}

type Context struct {
	Service  IService
	entityID EntityID

	pkt *codec.Packet

	err error

	// This mutex protects Keys map.
	mu sync.RWMutex
	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any
}

// func (c *Context) Copy() *Context {
// 	tc := &Context{
// 		entityID: c.entityID,
// 		pkt:      c.pkt,
// 		err:      c.err,
// 	}

// 	if c.pkt != nil {
// 		c.pkt.Retain()
// 	}

// 	return tc
// }

func (c *Context) reset() {
	// if c.pkt != nil {
	// 	c.pkt.Release()
	// }
	c.entityID = ""
	c.Service = nil
	c.pkt = nil
	c.err = nil
}

func (c *Context) Packet() *codec.Packet {
	return c.pkt
}

func (c *Context) EntityID() EntityID {
	return c.entityID
}

func (c *Context) Error(err error) {
	if err == nil {
		panic("err is nil")
	}
	c.err = err
	return
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline returns that there is no deadline (ok==false) when c.pkt has no Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.pkt == nil || c.pkt.Context() == nil {
		return
	}
	return c.pkt.Context().Deadline()
}

// Done returns nil (chan which will wait forever) when c.pkt has no Context.
func (c *Context) Done() <-chan struct{} {
	if c.pkt == nil || c.pkt.Context() == nil {
		return nil
	}
	return c.pkt.Context().Done()
}

// Err returns nil when c.pkt has no Context.
func (c *Context) Err() error {
	if c.pkt == nil || c.pkt.Context() == nil {
		return nil
	}
	return c.pkt.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key any) any {
	if key == 0 {
		return c.pkt
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	if c.pkt == nil || c.pkt.Context() == nil {
		return nil
	}
	return c.pkt.Context().Value(key)
}
