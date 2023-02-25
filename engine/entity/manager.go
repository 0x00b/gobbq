package entity

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[string]IService),
	entityDescs: make(map[string]*EntityDesc),
	Entities:    make(map[string]IBaseEntity),
}
var ProxyRegister RegisterProxy

type RegisterProxy interface {
	RegisterEntityToProxy(eid *bbq.EntityID) error
	RegisterServiceToProxy(svcName string) error

	// UnregisterEntityToProxy(eid EntityID) error
	// UnregisterServiceToProxy(svcName TypeName) error
}

// EntityManager manage entity lifecycle
type EntityManager struct {
	mu    sync.RWMutex // guards following
	serve bool

	Services    map[string]IService    // service name -> service info
	entityDescs map[string]*EntityDesc // entity name -> entity info
	Entities    map[string]IBaseEntity // entity id -> entity impl

}

func RegisterEntity(c Context, id *bbq.EntityID, entity IBaseEntity) error {
	ctx := allocContext(c)
	ctx.entity = entity
	entity.onInit(ctx, id)
	entity.OnInit()

	if c != nil {
		entity.setParant(c.Entity())
		c.Entity().addChildren(entity)
	}

	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	Manager.Entities[id.ID] = entity

	return nil
}

func NewEntity(c Context, id *bbq.EntityID) (IEntity, error) {
	desc, ok := Manager.entityDescs[id.Type]
	if !ok {
		xlog.Errorln("EntityManager.RegisterService new entity desc %s", id.Type)
		return nil, fmt.Errorf("EntityManager.RegisterService new entity desc %s", id.Type)
	}

	// new entity
	svcType := reflect.TypeOf(desc.EntityImpl)
	if svcType.Kind() == reflect.Pointer {
		svcType = svcType.Elem()
	}
	// 类型不对就在这里panic吧
	svcValue := reflect.New(svcType)
	svc := svcValue.Interface()
	entity, ok := svc.(IEntity)
	if !ok || entity == nil {
		xlog.Errorln("error type", svcType.Name())
		return nil, fmt.Errorf("new entity file %s", id.Type)
	}
	// init

	xlog.Debugln("register entity id:", id.String())

	newDesc := *desc
	newDesc.EntityImpl = entity
	entity.setDesc(&newDesc)

	RegisterEntity(c, id, entity)

	// start message loop
	go entity.Run()

	// send to poxy
	if ProxyRegister != nil {
		ProxyRegister.RegisterEntityToProxy(id)
	}

	return entity, nil
}

func (s *EntityManager) RegisterEntity(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Panicf("grpc: EntityManager.RegisterEntity found the handler of type %v that does not satisfy %v", st, ht)
			return
		}
	}
	s.registerEntity(sd, ss, intercepter...)
}

func (s *EntityManager) registerEntity(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
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
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	ss.setDesc(sd)

	s.registerServiceEntity(sd, ss)

	xlog.Tracef("grpc: EntityManager.RegisterService eid:%s", ss.EntityID())

	// start msg loop
	go ss.Run()

	if ProxyRegister != nil {
		ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
}

func (s *EntityManager) registerServiceEntity(sd *EntityDesc, entity IService) error {

	eid := entity.EntityID()
	if eid == nil || eid.ID == "" {
		if NewEntityID != nil {
			eid = NewEntityID.NewEntityID(sd.TypeName)
		} else {
			eid = &bbq.EntityID{ID: snowflake.GenUUID(), Type: sd.TypeName}
		}
	}

	RegisterEntity(nil, eid, entity)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	s.Services[sd.TypeName] = entity

	return nil
}
