package entity

import (
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/tool/snowflake"
	"github.com/0x00b/gobbq/xlog"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[TypeName]IEntity),
	entityDescs: make(map[TypeName]*EntityDesc),
	Entities:    make(map[EntityID]IEntity),
}
var ProxyRegister RegisterProxy

type RegisterProxy interface {
	RegisterEntityToProxy(eid EntityID) error
	RegisterServiceToProxy(svcName TypeName) error
}

// EntityManager manage entity lifecycle
type EntityManager struct {
	mu    sync.RWMutex // guards following
	serve bool

	Services    map[TypeName]IEntity     // service name -> service info
	entityDescs map[TypeName]*EntityDesc // entity name -> entity info
	Entities    map[EntityID]IEntity     // entity id -> entity impl

}

func NewEntity(c *Context, id *EntityID, typ TypeName) error {
	desc, ok := Manager.entityDescs[typ]
	if !ok {
		xlog.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", typ)
		return nil
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
	if !ok {
		xlog.Println("error type", svcType.Name())
		return nil
	}
	// init

	if id == nil || *id == "" {
		*id = EntityID(snowflake.GenUUID())
	}

	xlog.Println("register entity id:", *id)

	newDesc := *desc
	newDesc.EntityImpl = entity
	entity.setDesc(&newDesc)

	ctx := allocContext()
	ctx.Entity = entity
	entity.onInit(ctx, *id)

	if c != nil {
		entity.setParant(c.Entity)
		c.Entity.addChildren(entity)
	}

	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	Manager.Entities[*id] = entity

	// start message loop
	go entity.messageLoop()

	// send to poxy
	if ProxyRegister != nil {
		ProxyRegister.RegisterEntityToProxy(*id)
	}

	return nil
}

func (s *EntityManager) RegisterEntity(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Printf("grpc: EntityManager.RegisterEntity found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerEntity(sd, ss, intercepter...)
}

func (s *EntityManager) registerEntity(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	xlog.Printf("registerEntity(%q)", sd.TypeName)
	if s.serve {
		xlog.Printf("grpc: EntityManager.registerEntity after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.entityDescs[sd.TypeName]; ok {
		xlog.Printf("grpc: EntityManager.registerEntity found duplicate entity registration for %q", sd.TypeName)
		return
	}
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	s.entityDescs[sd.TypeName] = sd
}

func (s *EntityManager) RegisterService(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss, intercepter...)
}

func (s *EntityManager) registerService(sd *EntityDesc, ss IEntity, intercepter ...ServerInterceptor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	xlog.Printf("RegisterService(%q)", sd.TypeName)
	if s.serve {
		xlog.Printf("grpc: EntityManager.RegisterService after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.Services[sd.TypeName]; ok {
		xlog.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", sd.TypeName)
		return
	}
	sd.EntityImpl = ss
	sd.interceptors = intercepter
	ss.setDesc(sd)
	ctx := allocContext()
	ctx.Entity = ss
	ss.onInit(ctx, EntityID(snowflake.GenUUID()))

	s.Services[sd.TypeName] = ss

	// start msg loop
	go ss.messageLoop()

	if ProxyRegister != nil {
		ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
}
