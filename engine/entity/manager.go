package entity

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
)

var _ nets.PacketHandler = &EntityManager{}

// EntityManager manage entity lifecycle
type EntityManager struct {
	RemoteEntityManager

	mu    sync.RWMutex // guards following
	serve bool

	Services    map[string]IService    // service name -> service info
	entityDescs map[string]*EntityDesc // entity name -> entity info
	Entities    map[ID]IBaseEntity     // entity id -> entity impl

	ProxyRegister RegisterProxy

	EntityIDGenerator EntityIDGenerator
}

func NewEntityManager() *EntityManager {

	return &EntityManager{
		Services:    make(map[string]IService),
		entityDescs: make(map[string]*EntityDesc),
		Entities:    make(map[ID]IBaseEntity),
	}
}

type RemoteEntityManager interface {
	// for remote call, just send request packet, dont handle response
	SendPacket(pkt *codec.Packet) error
}

type RegisterProxy interface {
	// RegisterEntityToProxy(eid EntityID) error
	RegisterServiceToProxy(svcName string) error

	// UnregisterEntityToProxy(eid EntityID) error
	// UnregisterServiceToProxy(svcName TypeName) error
}

func (s *EntityManager) InitEntity(c Context, id EntityID, entity IBaseEntity) error {
	ctx, cancel := allocContext(c)
	ctx.entity = entity

	entity.onInit(ctx, cancel, id)
	entity.OnInit()

	if c != nil {
		entity.setParant(c.Entity())
		c.Entity().addChildren(entity)
	}
	return nil
}

func (s *EntityManager) RegisterEntity(c Context, id EntityID, entity IBaseEntity) error {

	s.InitEntity(c, id, entity)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entities[id.ID()] = entity

	return nil
}

func (s *EntityManager) ReplaceEntityID(c Context, old, new EntityID) error {

	s.mu.Lock()
	defer s.mu.Unlock()
	entity := s.Entities[old.ID()]
	s.Entities[new.ID()] = entity
	delete(s.Entities, old.ID())

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
	entity, ok := svc.(IEntity)
	if !ok || entity == nil {
		xlog.Errorln("error type", svcType.Name())
		return nil, fmt.Errorf("new entity file %s", typ)
	}
	// init

	xlog.Infoln("register entity id:", unsafe.Pointer(s), id.String())

	newDesc := *desc
	newDesc.EntityImpl = entity
	newDesc.EntityMgr = s
	entity.SetEntityDesc(&newDesc)

	s.RegisterEntity(c, id, entity)

	// start message loop
	wg := sync.WaitGroup{}
	wg.Add(1)
	go entity.Run(&wg)
	wg.Wait()

	// send to poxy
	// if s.ProxyRegister != nil {
	// s.ProxyRegister.RegisterEntityToProxy(id)
	// }

	return entity, nil
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
		v.OnDestroy()
		v.onDestroy()
	}

	// close entity
	for _, v := range s.Entities {
		v.OnDestroy()
		v.onDestroy()
	}

	return nil
}

func (s *EntityManager) registerEntityDesc(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {

	s.mu.Lock()
	defer s.mu.Unlock()
	xlog.Tracef("registerEntity(%q)", sd.TypeName)
	if s.serve {
		xlog.Tracef("gobbq: registerEntityDesc after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.entityDescs[sd.TypeName]; ok {
		xlog.Tracef("gobbq: registerEntityDesc found duplicate entity registration for %q", sd.TypeName)
		return
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
	if _, ok := s.Services[sd.TypeName]; ok {
		xlog.Tracef("gobbq: registerService found duplicate service registration for %q", sd.TypeName)
		return
	}
	sd.EntityMgr = s
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	ss.SetServiceDesc(sd)

	s.registerServiceEntity(sd, ss)

	xlog.Tracef("gobbq: registerService eid:%s", ss.EntityID())

	// start msg loop
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ss.Run(&wg)
	wg.Wait()
	if s.ProxyRegister != nil {
		s.ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		// s.ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
}

func (s *EntityManager) registerServiceEntity(sd *EntityDesc, entity IService) error {

	eid := entity.EntityID()
	if eid.Invalid() {
		if s.EntityIDGenerator == nil {
			return errors.New("no entity id generator")
		}
		eid = s.EntityIDGenerator.NewEntityID()
	}

	s.InitEntity(nil, eid, entity)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Services[sd.TypeName] = entity

	return nil
}
