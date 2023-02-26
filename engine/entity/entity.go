package entity

import (
	"errors"
	"sync"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

type EntityIDGenerator interface {
	NewEntityID(typeName string) *bbq.EntityID
}

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

	SetDesc(desc *EntityDesc)

	onInit(c Context, id *bbq.EntityID)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *codec.Packet)
	dispatchLocalCall(pkt *codec.Packet, req any, respChan chan any) error

	setParant(s IBaseEntity)
	addChildren(s IBaseEntity)
}

type IEntity interface {
	IBaseEntity

	entityType()
}

type methodHandler func(svc any, ctx Context, pkt *codec.Packet, interceptor ServerInterceptor)

type methodLocalHandler func(svc any, ctx Context, in any, interceptor ServerInterceptor) (any, error)

// MethodDesc represents an RPC Entity's method specification.
type MethodDesc struct {
	MethodName   string
	Handler      methodHandler
	LocalHandler methodLocalHandler
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

type localCall struct {
	pkt      *codec.Packet
	in       any
	respChan chan any
}

type baseEntity struct {
	// id
	entityID *bbq.EntityID

	// context
	context Context

	desc *EntityDesc

	callChan      chan *codec.Packet
	localCallChan chan *localCall

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
	xlog.Debugln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	defer func() {
		wg.Wait()

		xlog.Debugln("stop message loop", e.EntityID())
		// todo unregister entity

	}()

	// response
	go func() {
		for {
			select {
			case <-e.context.Done():
				xlog.Debugln("ctx done", e)
				return
			case pkt := <-e.respChan:
				wg.Add(1)
				xlog.Tracef("handle: %s", pkt.String())
				e.handleMethodRsp(e.context, pkt)
				wg.Done()
			}
		}
	}()

	// request, sync
	for {
		select {
		case <-e.context.Done():
			xlog.Debugln("ctx done", e)
			return

		case pkt := <-e.callChan:
			wg.Add(1)
			xlog.Tracef("handle: %s", pkt.String())
			err := e.handleCallMethod(e.context, pkt)
			if err != nil {
				xlog.Errorln(err)
			}
			xlog.Tracef("handle done: %s", pkt.String())
			wg.Done()
		case lc := <-e.localCallChan:
			wg.Add(1)
			xlog.Tracef("handle local call: %s", lc.pkt.String())
			err := e.handleLocalCallMethod(e.context, lc)
			if err != nil {
				xlog.Errorln(err)
			}
			xlog.Tracef("handle local call done: %s", lc.pkt.String())
			wg.Done()
		}
	}
}

//  for inner

func (e *baseEntity) registerCallback(requestID string, cb Callback) {
	if requestID == "" || cb == nil {
		return
	}

	xlog.Debugln("register callback:", requestID)

	e.cbMtx.Lock()
	defer e.cbMtx.Unlock()
	e.callback[requestID] = cb
}

func (e *baseEntity) onInit(c Context, id *bbq.EntityID) {
	e.context = c
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 1000)
	e.localCallChan = make(chan *localCall, 1000)
	e.callback = make(map[string]Callback, 1)
	e.respChan = make(chan *codec.Packet, 1)

	e.OnInit()
}

func (e *baseEntity) onDestroy() {
	e.OnDestroy()
	releaseContext(e.context)
}

func (e *baseEntity) SetDesc(desc *EntityDesc) {
	e.desc = desc
}

func (e *baseEntity) dispatchPkt(pkt *codec.Packet) {
	if pkt != nil {
		xlog.Traceln("dispatch:", e, pkt.String())
		pkt.Retain()
		if pkt.Header.RequestType == bbq.RequestType_RequestRequest {
			e.callChan <- pkt
			return
		}
		e.respChan <- pkt
	}
}

func (e *baseEntity) dispatchLocalCall(pkt *codec.Packet, req any, respChan chan any) error {
	if pkt != nil && req != nil {
		pkt.Retain()

		lc := localCall{
			pkt:      pkt,
			in:       req,
			respChan: respChan,
		}
		e.localCallChan <- &lc
		return nil
	}

	return ErrBadRequest
}

func (e *baseEntity) handleMethodRsp(c Context, pkt *codec.Packet) error {
	defer pkt.Release()

	c.setPacket(pkt)

	if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		cb, ok := e.callback[pkt.Header.RequestId]
		if ok {
			xlog.Debugln("callback:", pkt.Header.RequestId)
			e.cbMtx.Lock()
			defer e.cbMtx.Unlock()
			delete(e.callback, pkt.Header.RequestId)
			cb(pkt)
			return nil
		}
		xlog.Errorln("unknown response:", pkt.Header.RequestId)
		return errors.New("unknown response")
	}

	return nil
}

func (e *baseEntity) handleLocalCallMethod(c Context, lc *localCall) error {
	defer lc.pkt.Release()

	c.setPacket(lc.pkt)

	hdr := lc.pkt.Header

	sd := e.desc
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return ErrMethodNotFound
	}

	xlog.Traceln("LocalHandler 1", e.EntityID(), hdr.String())

	rsp, err := mt.LocalHandler(sd.EntityImpl, c, lc.in, chainServerInterceptors(sd.interceptors))

	xlog.Traceln("LocalHandler 2", hdr.String())

	if lc.respChan != nil {
		if rsp != nil {
			lc.respChan <- rsp
		} else {
			lc.respChan <- err
		}
	}

	return err
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
