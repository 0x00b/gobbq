package entity

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &EntityManager{}

// EntityManager manage entity lifecycle
type EntityManager struct {
	Proxy

	serve bool

	Services    map[string]IService    // service name -> service info
	entityDescs map[string]*EntityDesc // entity name -> entity info

	entityMtx sync.RWMutex         // guards following
	Entities  map[EntityID]IEntity // entity id -> entity impl

	ProxyRegister RegisterProxy

	EntityIDGenerator EntityIDGenerator
}

func NewEntityManager() *EntityManager {

	return &EntityManager{
		Services:    make(map[string]IService),
		entityDescs: make(map[string]*EntityDesc),
		Entities:    make(map[EntityID]IEntity),
	}
}

type Proxy interface {
	// for remote call, just send request packet, dont handle response
	SendPacket(pkt *nets.Packet) error
}

type RegisterProxy interface {
	// RegisterEntityToProxy(eid EntityID) error
	RegisterServiceToProxy(svcName string) error

	// UnregisterEntityToProxy(eid EntityID) error
	// UnregisterServiceToProxy(svcName TypeName) error
}

func (s *EntityManager) InitEntity(c Context, id EntityID, entity IEntity) error {
	ctx, cancel := allocContext(c)
	ctx.entity = entity

	entity.onInit(ctx, cancel, id)

	if c != nil {
		entity.setParant(c.Entity())
		c.Entity().addChildren(entity)
	}
	return nil
}

func (s *EntityManager) RegisterEntity(c Context, id EntityID, entity IEntity) error {

	s.InitEntity(c, id, entity)

	s.entityMtx.Lock()
	defer s.entityMtx.Unlock()
	s.Entities[id] = entity

	return nil
}

func (s *EntityManager) ReplaceEntityID(old, new EntityID) error {

	s.entityMtx.Lock()
	defer s.entityMtx.Unlock()
	entity := s.Entities[old]
	s.Entities[new] = entity
	delete(s.Entities, old)
	entity.setEntityID(new)

	return nil
}

func (s *EntityManager) NewEntity(c Context, id EntityID, typ string) (IEntity, error) {
	desc, ok := s.entityDescs[typ]
	if !ok {
		xlog.Errorln("NewEntity new entity desc ", typ, s, s.entityDescs)
		return nil, fmt.Errorf("NewEntity new entity desc %s", typ)
	}

	// new entity
	svcType := reflect.TypeOf(desc.EntityImpl)
	if svcType.Kind() == reflect.Pointer {
		svcType = svcType.Elem()
	}

	svcValue := reflect.New(svcType)
	svc := svcValue.Interface()
	e, ok := svc.(IEntity)
	if !ok || e == nil {
		xlog.Errorln("error type", svcType.Name())
		return nil, fmt.Errorf("new entity file %s", typ)
	}
	// init

	xlog.Infoln("register entity id:", id.String())

	newDesc := *desc
	newDesc.EntityImpl = e
	newDesc.EntityMgr = s
	SetEntityDesc(e, &newDesc)

	s.RegisterEntity(c, id, e)

	// start message loop
	Run(e)

	// send to poxy
	// if s.ProxyRegister != nil {
	// s.ProxyRegister.RegisterEntityToProxy(id)
	// }

	return e, nil
}

func (s *EntityManager) RegisterEntityDesc(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Panicf("gobbq: RegisterEntityDesc found the handler of type %v that does not satisfy %v", st, ht)
			return
		}
	}
	s.registerEntityDesc(sd, ss, intercepter...)
}

func (s *EntityManager) Close(ch chan struct{}) error {
	// close svc
	for _, v := range s.Services {
		v.Stop()
	}

	// close entity
	func() {
		s.entityMtx.Lock()
		defer s.entityMtx.Unlock()

		for _, v := range s.Entities {
			v.Stop()
		}

	}()
	return nil
}

func (s *EntityManager) registerEntityDesc(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {

	xlog.Tracef("registerEntity(%q)", sd.TypeName)
	if s.serve {
		xlog.Tracef("gobbq: registerEntityDesc after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.entityDescs[sd.TypeName]; ok {
		xlog.Tracef("gobbq: registerEntityDesc found duplicate entity registration for %q", sd.TypeName)
		return
	}

	for k, v := range BbqSysEntityDesc.Methods {
		if _, ok := sd.Methods[k]; ok {
			panic("dup method with sys:" + k)
		}
		sd.Methods[k] = v
	}

	sd.EntityMgr = s
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	s.entityDescs[sd.TypeName] = sd

	xlog.Traceln("registerEntityDesc", sd)
}

func (s *EntityManager) RegisterService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Tracef("gobbq: RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss, intercepter...)
}

func (s *EntityManager) registerService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
	xlog.Tracef("RegisterService(%q)", sd.TypeName)
	if s.serve {
		xlog.Tracef("gobbq: registerService after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.GetService(sd.TypeName); ok {
		xlog.Tracef("gobbq: registerService found duplicate service registration for %q", sd.TypeName)
		return
	}

	for k, v := range BbqSysEntityDesc.Methods {
		if _, ok := sd.Methods[k]; ok {
			panic("dup method with sys:" + k)
		}
		sd.Methods[k] = v
	}

	sd.EntityMgr = s
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	SetServiceDesc(ss, sd)

	xlog.Tracef("gobbq: registerService 111 eid:%d", ss.EntityID())

	s.registerServiceEntity(sd, ss)

	xlog.Tracef("gobbq: registerService 222 eid:%d", ss.EntityID())

	// start msg loop
	Run(ss)

	xlog.Tracef("gobbq: registerService 333 eid:%d", ss.EntityID())

	if s.ProxyRegister != nil {
		s.ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		// s.ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
	xlog.Tracef("gobbq: registerService xxxxxx eid:%d", ss.EntityID())

}

func (s *EntityManager) registerServiceEntity(sd *EntityDesc, entity IService) error {

	xlog.Tracef("gobbq: registerService 444 eid")

	// eid := entity.EntityID()
	// if eid.Invalid() {
	// 	if s.EntityIDGenerator == nil {
	// 		return errors.New("no entity id generator")
	// 	}
	eid := s.EntityIDGenerator.NewEntityID()
	// }

	xlog.Tracef("gobbq: registerService 555 eid")
	s.RegisterEntity(nil, eid, entity)
	xlog.Tracef("gobbq: registerService 666 eid")

	s.Services[sd.TypeName] = entity

	xlog.Tracef("gobbq: registerService 777 eid")

	return nil
}

// 需要优化, gate的id再拆分,不要直接用这个来判断是不是自己的entity
func (s *EntityManager) GetEntity(eid EntityID) (IEntity, bool) {
	s.entityMtx.RLock()
	defer s.entityMtx.RUnlock()

	e, ok := s.Entities[eid]
	return e, ok
}

// 需要优化, gate的id再拆分,不要直接用这个来判断是不是自己的entity
func (s *EntityManager) removeEntity(eid EntityID) bool {
	s.entityMtx.Lock()
	defer s.entityMtx.Unlock()

	delete(s.Entities, eid)

	return true
}

func (s *EntityManager) GetService(typ string) (IService, bool) {
	svc, ok := s.Services[typ]
	return svc, ok
}
