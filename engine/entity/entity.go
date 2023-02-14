package entity

import (
	"errors"
	"sync"

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

	Run()

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in

	// for inner

	setDesc(desc *EntityDesc)

	onInit(c *Context, id EntityID)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *codec.Packet)

	setParant(s IEntity)
	addChildren(s IEntity)
}

type methodHandler func(svc any, ctx *Context, pkt *codec.Packet, interceptor ServerInterceptor)

// type methodLocalHandler func(svc any, ctx *Context, in any, interceptor ServerInterceptor) (any, error)

// MethodDesc represents an RPC Entity's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
	// LocalHandler methodLocalHandler
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

	callChan chan *codec.Packet
	respChan chan *codec.Packet

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

func (e *Entity) Context() *Context {
	return e.context
}

func (e *Entity) onInit(c *Context, id EntityID) {
	e.context = c
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 1000)
	e.callback = make(map[string]Callback, 1)
	e.respChan = make(chan *codec.Packet, 1)

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
		xlog.Println("dispatch:", e, pkt.String())
		pkt.Retain()
		if pkt.Header.RequestType == bbq.RequestType_RequestRequest {
			e.callChan <- pkt
			return
		}
		e.respChan <- pkt
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

func (e *Entity) Run() {
	xlog.Println("start message loop", e.EntityID())

	go func() {
		for !e.done {
			pkt := <-e.respChan
			xlog.Printf("handle: %s", pkt.String())
			e.handleMethodRsp(e.context, pkt)
		}
	}()

	for !e.done {
		select {
		case pkt := <-e.callChan:
			xlog.Printf("handle: %s", pkt.String())
			e.handleCallMethod(e.context, pkt)
		}
	}
	xlog.Println("stop message loop", e.EntityID())
	// todo unregister entity
}

func (e *Entity) handleMethodRsp(c *Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.pkt = pkt

	if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		cb, ok := e.callback[pkt.Header.RequestId]
		if ok {
			xlog.Println("callback:", pkt.Header.RequestId)
			e.cbMtx.Lock()
			defer e.cbMtx.Unlock()
			delete(e.callback, pkt.Header.RequestId)
			cb(pkt)
			return nil
		}
		xlog.Println("unknown response:", pkt.Header.RequestId)
		return errors.New("unknown response")
	}

	return nil
}

func (e *Entity) handleCallMethod(c *Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.pkt = pkt

	sd := e.desc
	// todo method name repeat get
	hdr := pkt.Header
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return ErrMethodNotFound
	}

	mt.Handler(sd.EntityImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

func (e *Entity) setParant(svc IEntity) {

}

func (e *Entity) addChildren(ety IEntity) {

}

func (e *Entity) EntityID() EntityID {
	return e.entityID
}
