package entity

import (
	"errors"
	"sync"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

// just for inner
// type EntityID bbq.EntityID

type EntityIDGenerator interface {
	NewEntityID(typeName string) *bbq.EntityID
}

var NewEntityID EntityIDGenerator

// func ToEntityID(id *bbq.EntityID) *EntityID {
// 	if id == nil {
// 		return nil
// 	}
// 	return &EntityID{
// 		ID:      id.ID,
// 		Type:    TypeName(id.Type),
// 		ProxyID: id.ProxyID,
// 	}
// }

// func ToPBEntityID(id *EntityID) *bbq.EntityID {
// 	if id == nil {
// 		return nil
// 	}

// 	return &bbq.EntityID{
// 		ID:      id.ID,
// 		Type:    string(id.Type),
// 		ProxyID: id.ProxyID,
// 	}
// }

// just for inner
// type TypeName string

var _ IEntity = &Entity{}

type Entity struct {
	baseEntity
}

func (e *Entity) entityType() {}

type IBaseEntity interface {

	// EntityID
	EntityID() *bbq.EntityID

	// Entity Lifetime
	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnDestroy() // Called when entity is destroying (just before destroy)

	Desc() *EntityDesc

	Context() Context

	Run()

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in

	// Watch/unwatch entity

	// for inner

	registerCallback(requestID string, cb Callback)

	setDesc(desc *EntityDesc)

	onInit(c Context, id *bbq.EntityID)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *codec.Packet)

	setParant(s IBaseEntity)
	addChildren(s IBaseEntity)
}

type IEntity interface {
	IBaseEntity

	entityType()
}

type methodHandler func(svc any, ctx Context, pkt *codec.Packet, interceptor ServerInterceptor)

// type methodLocalHandler func(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error)

// MethodDesc represents an RPC Entity's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
	// LocalHandler methodLocalHandler
}

// EntityDesc represents an RPC Entity's specification.
type EntityDesc struct {
	EntityImpl any

	TypeName string
	// The pointer to the Entity interface. Used to check whether the user
	// provided implementation satisfies the interface requiremente.
	HandlerType any
	Methods     map[string]MethodDesc
	Metadata    any

	interceptors []ServerInterceptor
}

type baseEntity struct {
	// id
	entityID *bbq.EntityID

	// context
	context Context

	desc *EntityDesc

	callChan chan *codec.Packet
	respChan chan *codec.Packet

	cbMtx sync.RWMutex
	// requestid -> callback
	callback map[string]Callback
}

func (e *baseEntity) OnInit() {}

func (*baseEntity) OnDestroy() {}

func (e *baseEntity) Desc() *EntityDesc {
	return e.desc
}

// Migration
func (e *baseEntity) OnMigrateOut() {} // Called just before entity is migrating out
func (e *baseEntity) OnMigrateIn()  {} // Called just after entity is migrating in

func (e *baseEntity) Context() Context {
	return e.context
}

func (e *baseEntity) EntityID() *bbq.EntityID {
	return e.entityID
}

func (e *baseEntity) Run() {
	xlog.Println("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	defer func() {
		wg.Wait()

		xlog.Println("stop message loop", e.EntityID())
		// todo unregister entity

	}()

	wg.Add(1)

	// response
	go func() {
		defer wg.Done()
		for {
			select {
			case <-e.context.Done():
				xlog.Println("ctx done", e)
			case pkt := <-e.respChan:
				xlog.Printf("handle: %s", pkt.String())
				e.handleMethodRsp(e.context, pkt)
			}
		}
	}()

	// request, sync
	for {
		select {
		case <-e.context.Done():
			xlog.Println("ctx done", e)

		case pkt := <-e.callChan:
			xlog.Printf("handle: %s", pkt.String())
			e.handleCallMethod(e.context, pkt)
		}
	}
}

//  for inner

func (e *baseEntity) registerCallback(requestID string, cb Callback) {
	if requestID == "" || cb == nil {
		return
	}
	e.cbMtx.Lock()
	defer e.cbMtx.Unlock()
	e.callback[requestID] = cb
}

func (e *baseEntity) onInit(c Context, id *bbq.EntityID) {
	e.context = c
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 1000)
	e.callback = make(map[string]Callback, 1)
	e.respChan = make(chan *codec.Packet, 1)

	e.OnInit()
}

func (e *baseEntity) onDestroy() {
	e.OnDestroy()
	releaseContext(e.context)
}

func (e *baseEntity) setDesc(desc *EntityDesc) {
	e.desc = desc
}

func (e *baseEntity) dispatchPkt(pkt *codec.Packet) {
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

func (e *baseEntity) handleMethodRsp(c Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.setPacket(pkt)

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

func (e *baseEntity) handleCallMethod(c Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.setPacket(pkt)
	sd := e.desc

	hdr := pkt.Header
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return ErrMethodNotFound
	}

	mt.Handler(sd.EntityImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

func (e *baseEntity) setParant(svc IBaseEntity) {
}

func (e *baseEntity) addChildren(ety IBaseEntity) {
}
