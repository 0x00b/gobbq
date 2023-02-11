package entity

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

// just for inner
type EntityID string

// just for inner
type TypeName string

type IEntity interface {

	// EntityID
	EntityID() EntityID

	// Entity Lifetime
	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnDestroy() // Called when entity is destroying (just before destroy)

	Desc() *EntityDesc

	RegisterCallback(requestID string, cb Callback)

	Context() *Context

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in
	// Freeze && Restore
	OnFreeze()   // Called when entity is freezing
	OnRestored() // Called when entity is restored

	// for inner

	setDesc(desc *EntityDesc)

	onInit(c *Context, id EntityID)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *codec.Packet)
	messageLoop()

	setParant(s IEntity)
	addChildren(s IEntity)
}

type methodHandler func(svc any, ctx *Context, pkt *codec.Packet, interceptor ServerInterceptor)
type methodLocalHandler func(svc any, ctx *Context, in any, callback func(c *Context, rsp any), interceptor ServerInterceptor)

// MethodDesc represents an RPC Entity's method specification.
type MethodDesc struct {
	MethodName   string
	Handler      methodHandler
	LocalHandler methodLocalHandler
}

// EntityDesc represents an RPC Entity's specification.
type EntityDesc struct {
	EntityImpl any

	TypeName TypeName
	// The pointer to the Entity interface. Used to check whether the user
	// provided implementation satisfies the interface requiremente.
	HandlerType any
	Methods     map[string]MethodDesc
	Metadata    any

	interceptors []ServerInterceptor
}

var _ IEntity = &Entity{}

type Entity struct {
	// id
	entityID EntityID

	// context
	context *Context

	desc *EntityDesc

	done bool

	pktChan chan *codec.Packet

	cbMtx sync.RWMutex
	// requestid -> callback
	callback map[string]Callback
}

func (e *Entity) OnInit() {
}

func (*Entity) OnDestroy() {
	// Called when entity is destroying (just before destroy)
}

func (e *Entity) Desc() *EntityDesc {
	return e.desc
}

// Migration
func (e *Entity) OnMigrateOut() {} // Called just before entity is migrating out
func (e *Entity) OnMigrateIn()  {} // Called just after entity is migrating in
// Freeze && Restore
func (e *Entity) OnFreeze()   {} // Called when entity is freezing
func (e *Entity) OnRestored() {} // Called when entity is restored

func (e *Entity) Context() *Context {
	return e.context
}

func (e *Entity) onInit(c *Context, id EntityID) {
	e.context = c
	e.entityID = id

	e.callback = make(map[string]Callback, 10000)
	e.pktChan = make(chan *codec.Packet, 10000)

	e.OnInit()
}

func (e *Entity) onDestroy() {
	e.done = true
	e.OnDestroy()
	releaseContext(e.context)
}

func (e *Entity) setDesc(desc *EntityDesc) {
	e.desc = desc
}

func (e *Entity) dispatchPkt(pkt *codec.Packet) {
	if pkt != nil {
		xlog.Printf("dispatch:%d %d %s", unsafe.Pointer(pkt), pkt, pkt.String())
		pkt.Retain()
		e.pktChan <- pkt
	}
}
func (e *Entity) RegisterCallback(requestID string, cb Callback) {
	if requestID == "" || cb == nil {
		return
	}
	e.cbMtx.Lock()
	defer e.cbMtx.Unlock()
	e.callback[requestID] = cb
}

func (e *Entity) messageLoop() {
	xlog.Println("start message loop")
	for !e.done {
		select {
		case pkt := <-e.pktChan:
			xlog.Printf("handle:%d %s", pkt, pkt.String())
			e.handleCallMethod(e.context, pkt)
		}
	}
}

func (e *Entity) handleCallMethod(c *Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.pkt = pkt

	if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		cb, ok := e.callback[pkt.Header.RequestId]
		if ok {
			cb(c, pkt)
			return nil
		}
		return errors.New("unknown response")
	}

	sd := e.desc
	// todo method name repeat get
	hdr := pkt.Header
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return MethodNotFound
	}

	mt.Handler(sd.EntityImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

func (e *Entity) setContext(c *Context) {
	e.context = c
}

func (e *Entity) setParant(svc IEntity) {

}

func (e *Entity) addChildren(ety IEntity) {

}

func (e *Entity) EntityID() EntityID {
	return e.entityID
}

func (e *Entity) setEntityID(id EntityID) {
	e.entityID = id
	return
}
