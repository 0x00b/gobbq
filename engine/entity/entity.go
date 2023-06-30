package entity

import (
	"errors"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/engine/timer"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
)

var _ IEntity = &Entity{}

// On开头的表示可重写
type IEntity interface {

	// EntityID
	EntityID() EntityID
	Context() Context
	EntityDesc() *EntityDesc

	// Entity Lifetime
	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnDestroy() // Called when entity is destroying (just before destroy)

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in

	// AddCallback 直接用，不需要重写，除非有特殊需求
	AddCallback(d time.Duration, callback timer.CallbackFunc)
	// AddTimer 直接用，不需要重写，除非有特殊需求
	AddTimer(d time.Duration, callback timer.CallbackFunc)
	// OnTick 实时性很高要求的可以通过tick实现，service最低tick时间5ms， entity执行过一次事件之后执行一次OnTick
	OnTick()

	// 关注entity, 如果entity退出或者状态变更会通过Receive接收到状态变更通知
	Watch(id EntityID)
	Unwatch(id EntityID)
	OnNotify(NotifyInfo)

	// 主动结束
	Stop()

	// === for inner ===
	getEntityMgr() *EntityManager

	registerCallback(requestID string, cb Callback)
	popCallback(requestID string) (Callback, bool)

	onInit(c Context, cancel func(), id EntityID)
	setEntityID(id EntityID)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *nets.Packet)
	dispatchLocalCall(pkt *nets.Packet, req any, respChan chan any) error

	setParant(s IEntity)
	addChildren(s IEntity)

	setEntityDesc(desc *EntityDesc)

	entityType()

	run(ch chan bool)
}

type State int

type NotifyInfo struct {
	EntityID EntityID
}

type methodHandler func(svc any, ctx Context, pkt *nets.Packet, interceptor ServerInterceptor)

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

	EntityMgr *EntityManager
}

type localCall struct {
	pkt      *nets.Packet
	in       any
	respChan chan any
}

type Entity struct {
	// id
	entityID EntityID

	// 	// status
	// （1）0，待运行：异常重启后恢复，未恢复完前的状态
	// （2）1，正常状态：处于这个状态时，可以被分配用作帧消息转发
	// （3）2，迁移：下线前，需要将对局都迁移到新的，迁移状态不接收新的对局
	// （4）3，异常：的心跳丢失，异常状态持续一段时间后，进入迁移状态
	// （5）4，销毁：完成迁移后，会进入销毁状态

	desc *EntityDesc

	// context
	context Context
	cancel  func()

	callChan      chan *nets.Packet
	localCallChan chan *localCall

	respChan chan *nets.Packet

	cbMtx sync.RWMutex
	// requestid -> callback
	callback map[string]Callback

	timer timer.Timer

	ticker *time.Ticker

	runOnce     sync.Once
	initOnce    sync.Once
	destroyOnce sync.Once

	watchers map[EntityID]bool
}

func Run(e IEntity) {
	ch := make(chan bool)
	secure.GO(func() {
		e.run(ch)
	})
	<-ch
}

func (e *Entity) run(ch chan bool) {
	done := true
	e.runOnce.Do(func() {
		done = false
	})
	if done {
		ch <- true
		close(ch)
		return
	}

	xlog.Debugln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	tempch := make(chan bool)
	defer func() {
		wg.Wait()
		xlog.Debugln("stop message loop", e.EntityID())

		// todo unregister entity

		e.onDestroy()

		close(tempch)
		close(ch)
	}()

	// response
	secure.GO(func() {
		for {
			select {
			case tempch <- true:

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
	})

	// 上面的for执行了，在继续下面的for
	<-tempch

	// request, sync
	for {
		select {
		case ch <- true:

		case <-e.context.Done():
			xlog.Debugln("ctx done", e)
			return

		case pkt := <-e.callChan:
			wg.Add(1)
			xlog.Tracef("handle: %s", pkt.String())
			err := e.handleCallMethod(e.context, pkt, e.EntityDesc())
			if err != nil {
				xlog.Errorln(err)
			}
			xlog.Tracef("handle done: %s", pkt.String())
			wg.Done()

		case lc := <-e.localCallChan:
			wg.Add(1)
			xlog.Tracef("handle local call: %s", lc.pkt.String())
			err := e.handleLocalCallMethod(e.context, lc, e.EntityDesc())
			if err != nil {
				xlog.Errorln(err)
			}
			xlog.Tracef("handle local call done: %s", lc.pkt.String())
			wg.Done()

		case <-e.ticker.C:
			e.timer.Tick()
		}

		// 实时性很高要求的可以通过tick实现
		e.context.Entity().OnTick()
	}
}

func (e *Entity) Stop() {
	e.cancel()
	e.getEntityMgr().removeEntity(e.EntityID())
}

//  for inner

func (e *Entity) entityType() {}

func (e *Entity) EntityDesc() *EntityDesc {
	return e.desc
}

func SetEntityDesc(e IEntity, desc *EntityDesc) {
	e.setEntityDesc(desc)
}

func (e *Entity) setEntityDesc(desc *EntityDesc) {
	e.desc = desc
}

func (e *Entity) getEntityMgr() *EntityManager {
	return e.desc.EntityMgr
}

func (e *Entity) registerCallback(requestID string, cb Callback) {
	if requestID == "" || cb == nil {
		return
	}

	// xlog.Debugln("register callback:", requestID)

	e.cbMtx.Lock()
	defer e.cbMtx.Unlock()
	e.callback[requestID] = cb
}

func (e *Entity) popCallback(requestID string) (Callback, bool) {
	if requestID == "" {
		return nil, false
	}

	// xlog.Debugln("unregister callback:", requestID)

	e.cbMtx.Lock()
	defer e.cbMtx.Unlock()
	cb, ok := e.callback[requestID]
	if !ok {
		return nil, false
	}
	delete(e.callback, requestID)

	return cb, true
}

func (e *Entity) setEntityID(id EntityID) {
	e.entityID = id
}

func (e *Entity) onInit(c Context, cancel func(), id EntityID) {
	e.initOnce.Do(func() {
		e.setEntityID(id)

		e.context = c
		e.cancel = cancel
		e.callChan = make(chan *nets.Packet, 1000)
		e.localCallChan = make(chan *localCall, 1000)
		e.callback = make(map[string]Callback, 1)
		e.respChan = make(chan *nets.Packet, 1)
		e.timer.Init()
		e.ticker = time.NewTicker(GAME_SERVICE_TICK_INTERVAL)
		e.watchers = make(map[EntityID]bool)

		c.Entity().OnInit()
	})
}

func (e *Entity) onDestroy() {
	e.destroyOnce.Do(func() {

		e.context.Entity().OnDestroy()

		for id := range e.watchers {
			client := NewBbqSysClient(id)
			client.SysNotify(e.Context(), &WatchRequest{EntityID: uint64(e.EntityID())})
		}

		e.ticker.Stop()

		close(e.callChan)
		close(e.respChan)
		close(e.localCallChan)

		releaseContext(e.context)
	})
}

func DispatchPkt(e IEntity, pkt *nets.Packet) {
	e.dispatchPkt(pkt)
}

func (e *Entity) dispatchPkt(pkt *nets.Packet) {
	if pkt != nil {
		// xlog.Traceln("dispatch:", pkt.String())
		pkt.Retain()
		if pkt.Header.RequestType == bbq.RequestType_RequestRequest {
			e.callChan <- pkt
			return
		}
		e.respChan <- pkt
	}
}

func (e *Entity) dispatchLocalCall(pkt *nets.Packet, req any, respChan chan any) error {
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

func (e *Entity) initContext(c Context, pkt *nets.Packet) {
	c.setPacket(pkt)
	// SetEntityMgr(c, e.desc.EntityMgr)
}

func (e *Entity) handleMethodRsp(c Context, pkt *nets.Packet) error {
	defer pkt.Release()

	e.initContext(c, pkt)

	if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		cb, ok := e.popCallback(pkt.Header.RequestId)
		if ok {
			// xlog.Debugln("callback:", pkt.Header.RequestId)
			cb(pkt)
			return nil
		}
		xlog.Errorln("unknown response:", pkt.Header.RequestId)
		return errors.New("unknown response")
	}

	return nil
}

func (e *Entity) handleLocalCallMethod(c Context, lc *localCall, sd *EntityDesc) error {
	defer lc.pkt.Release()

	e.initContext(c, lc.pkt)

	hdr := lc.pkt.Header

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

func (e *Entity) handleCallMethod(c Context, pkt *nets.Packet, sd *EntityDesc) error {
	defer pkt.Release()

	e.initContext(c, pkt)

	hdr := pkt.Header
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return ErrMethodNotFound
	}

	mt.Handler(sd.EntityImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

// AddOnceTimer 直接用，不需要重写，除非有特殊需求
func (e *Entity) AddCallback(d time.Duration, callback timer.CallbackFunc) {
	e.timer.AddCallback(d, callback)
}

// AddRepeatTimer 直接用，不需要重写，除非有特殊需求
func (e *Entity) AddTimer(d time.Duration, callback timer.CallbackFunc) {
	e.timer.AddTimer(d, callback)
}

// empty/default implement

func (e *Entity) setParant(svc IEntity)   {}
func (e *Entity) addChildren(ety IEntity) {}

func (e *Entity) OnInit()            {}
func (e *Entity) OnDestroy()         {}
func (e *Entity) OnMigrateOut()      {} // Called just before entity is migrating out
func (e *Entity) OnMigrateIn()       {} // Called just after entity is migrating in
func (e *Entity) Context() Context   { return e.context }
func (e *Entity) OnTick()            {}
func (e *Entity) EntityID() EntityID { return e.entityID }
