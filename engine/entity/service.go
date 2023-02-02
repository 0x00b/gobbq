package entity

import (
	"fmt"
	"sync"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
)

type IService interface {
	// Entity Lifetime
	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnDestroy() // Called when entity is destroying (just before destroy)

	Desc() *ServiceDesc

	RegisterCallback(requestID string, cb Callback)

	Context() *Context

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called just after entity is migrating in
	// Freeze && Restore
	OnFreeze()   // Called when entity is freezing
	OnRestored() // Called when entity is restored

	// for inner

	setDesc(desc *ServiceDesc)

	onInit(*Context)
	onDestroy() // Called when entity is destroying (just before destroy), for inner

	dispatchPkt(pkt *codec.Packet)
	messageLoop()

	setParant(s IService)
	addChildren(s IService)
}

type methodHandler func(svc any, ctx *Context, pkt *codec.Packet, interceptor ServerInterceptor)
type methodLocalHandler func(svc any, ctx *Context, in any, callback func(c *Context, rsp any), interceptor ServerInterceptor)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName   string
	Handler      methodHandler
	LocalHandler methodLocalHandler
}

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct {
	ServiceImpl any

	TypeName TypeName
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType any
	Methods     map[string]MethodDesc
	Metadata    any

	interceptors []ServerInterceptor
}

var _ IService = &Service{}

type Service struct {
	context *Context

	desc *ServiceDesc

	done bool

	pktChan chan *codec.Packet

	cbMtx sync.RWMutex
	// requestid -> callback
	callback map[string]Callback
}

func (*Service) OnInit() {
	// Called when initializing entity struct, override to initialize entity custom fields
}

func (*Service) OnDestroy() {
	// Called when entity is destroying (just before destroy)
}

func (s *Service) Desc() *ServiceDesc {
	return s.desc
}

// Migration
func (s *Service) OnMigrateOut() {} // Called just before entity is migrating out
func (s *Service) OnMigrateIn()  {} // Called just after entity is migrating in
// Freeze && Restore
func (s *Service) OnFreeze()   {} // Called when entity is freezing
func (s *Service) OnRestored() {} // Called when entity is restored

func (s *Service) Context() *Context {
	return s.context
}

func (s *Service) onInit(c *Context) {
	s.callback = make(map[string]Callback)
	s.pktChan = make(chan *codec.Packet, 10000)
	s.context = c

	s.OnInit()
}

func (s *Service) onDestroy() {
	s.done = true
	s.OnDestroy()
	releaseContext(s.context)
}

func (s *Service) setDesc(desc *ServiceDesc) {
	s.desc = desc
}

func (s *Service) dispatchPkt(pkt *codec.Packet) {
	if pkt != nil {
		s.pktChan <- pkt
	}
}
func (s *Service) RegisterCallback(requestID string, cb Callback) {
	if requestID == "" || cb == nil {
		return
	}
	s.cbMtx.Lock()
	defer s.cbMtx.Unlock()
	s.callback[requestID] = cb
}

func (s *Service) messageLoop() {
	fmt.Println("start message loop")
	for !s.done {
		select {
		case pkt := <-s.pktChan:
			s.handleCallMethod(s.context, pkt)
		}
	}
}

func (s *Service) handleCallMethod(c *Context, pkt *codec.Packet) error {
	c.pkt = pkt

	if pkt.Header.RequestType == bbq.RequestType_RequestRespone {
		cb, ok := s.callback[pkt.Header.RequestId]
		if ok {
			cb(c, nil)
			return nil
		}
	}

	sd := s.desc
	// todo method name repeat get
	hdr := pkt.Header
	mt, ok := sd.Methods[hdr.Method]
	if !ok {
		return MethodNotFound
	}

	mt.Handler(sd.ServiceImpl, c, pkt, chainServerInterceptors(sd.interceptors))

	return nil
}

func (s *Service) setContext(c *Context) {
	s.context = c
}

func (s *Service) setParant(svc IService) {

}

func (s *Service) addChildren(ety IService) {

}
