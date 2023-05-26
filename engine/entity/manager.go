package entity

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
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
	Entities    map[string]IBaseEntity // entity id -> entity impl

	ProxyRegister RegisterProxy

	EntityIDGenerator EntityIDGenerator
}

func NewEntityManager() *EntityManager {

	return &EntityManager{
		Services:    make(map[string]IService),
		entityDescs: make(map[string]*EntityDesc),
		Entities:    make(map[string]IBaseEntity),
	}
}

type RemoteEntityManager interface {
	// for remote call, just send request packet, dont handle response
	SendPackt(pkt *codec.Packet) error
}

type RegisterProxy interface {
	// RegisterEntityToProxy(eid *bbq.EntityID) error
	RegisterServiceToProxy(svcName string) error

	// UnregisterEntityToProxy(eid EntityID) error
	// UnregisterServiceToProxy(svcName TypeName) error
}

func (s *EntityManager) RegisterEntity(c Context, id *bbq.EntityID, entity IBaseEntity) error {
	ctx, cancel := allocContext(c)
	ctx.entity = entity

	entity.onInit(ctx, cancel, id)
	entity.OnInit()

	if c != nil {
		entity.setParant(c.Entity())
		c.Entity().addChildren(entity)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entities[id.ID] = entity

	return nil
}

func (s *EntityManager) NewEntity(c Context, id *bbq.EntityID) (IEntity, error) {
	desc, ok := s.entityDescs[id.Type]
	if !ok {
		xlog.Errorln("EntityManager.RegisterService new entity desc %s", id.Type)
		return nil, fmt.Errorf("EntityManager.RegisterService new entity desc %s", id.Type)
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
		return nil, fmt.Errorf("new entity file %s", id.Type)
	}
	// init

	xlog.Infoln("register entity id:", unsafe.Pointer(s), id.String())

	newDesc := *desc
	newDesc.EntityImpl = entity
	newDesc.EntityMgr = s
	entity.SetDesc(&newDesc)

	s.RegisterEntity(c, id, entity)

	// start message loop
	go entity.Run()

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
			xlog.Panicf("grpc: EntityManager.RegisterEntity found the handler of type %v that does not satisfy %v", st, ht)
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
		xlog.Tracef("grpc: EntityManager.registerEntity after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.entityDescs[sd.TypeName]; ok {
		xlog.Tracef("grpc: EntityManager.registerEntity found duplicate entity registration for %q", sd.TypeName)
		return
	}
	sd.EntityMgr = s
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	s.entityDescs[sd.TypeName] = sd
}

func (s *EntityManager) RegisterService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Tracef("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss, intercepter...)
}

func (s *EntityManager) registerService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
	xlog.Tracef("RegisterService(%q)", sd.TypeName)
	if s.serve {
		xlog.Tracef("grpc: EntityManager.RegisterService after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.Services[sd.TypeName]; ok {
		xlog.Tracef("grpc: EntityManager.RegisterService found duplicate service registration for %q", sd.TypeName)
		return
	}
	sd.EntityMgr = s
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	ss.SetDesc(sd)

	s.registerServiceEntity(sd, ss)

	xlog.Tracef("grpc: EntityManager.RegisterService eid:%s", ss.EntityID())

	// start msg loop
	go ss.Run()

	if s.ProxyRegister != nil {
		s.ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		// s.ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
}

func (s *EntityManager) registerServiceEntity(sd *EntityDesc, entity IService) error {

	eid := entity.EntityID()
	if eid == nil || eid.ID == "" {
		if s.EntityIDGenerator != nil {
			eid = s.EntityIDGenerator.NewEntityID(sd.TypeName)
		} else {
			eid = &bbq.EntityID{ID: snowflake.GenUUID(), Type: sd.TypeName}
		}
	}

	s.RegisterEntity(nil, eid, entity)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Services[sd.TypeName] = entity

	return nil
}
