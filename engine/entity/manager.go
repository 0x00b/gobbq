package entity

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/snowflake"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[TypeName]IService),
	entityDescs: make(map[TypeName]*ServiceDesc),
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

	Services    map[TypeName]IService     // service name -> service info
	entityDescs map[TypeName]*ServiceDesc // entity name -> entity info
	Entities    map[EntityID]IEntity      // entity id -> entity impl

}

func NewEntity(c *Context, id EntityID, typ TypeName) *bbq.EntityID {
	desc, ok := Manager.entityDescs[typ]
	if !ok {
		fmt.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", typ)
		return nil
	}

	// new entity
	svcType := reflect.TypeOf(desc.ServiceImpl)
	if svcType.Kind() == reflect.Pointer {
		svcType = svcType.Elem()
	}
	// 类型不对就在这里panic吧
	svcValue := reflect.New(svcType)
	svc := svcValue.Interface()
	entity, ok := svc.(IEntity)
	if !ok {
		fmt.Println("error type", svcType.Name())
		return nil
	}
	// init

	if id == "" {
		id = EntityID(snowflake.GenUUID())
	}

	fmt.Println("register entity id:", id)

	newDesc := *desc
	newDesc.ServiceImpl = entity
	entity.setDesc(&newDesc)

	ctx := allocContext()
	ctx.Service = entity
	ctx.entityID = id
	entity.onInit(ctx)

	if c != nil {
		entity.setParant(c.Service)
		c.Service.addChildren(entity)
	}

	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	Manager.Entities[id] = entity

	// start message loop
	go entity.messageLoop()

	// send to poxy
	if ProxyRegister != nil {
		ProxyRegister.RegisterEntityToProxy(id)
	}

	return &bbq.EntityID{ID: string(id), TypeName: string(desc.TypeName)}
}

func (s *EntityManager) RegisterEntity(sd *ServiceDesc, ss IEntity, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			fmt.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerEntity(sd, ss, intercepter...)
}

func (s *EntityManager) registerEntity(sd *ServiceDesc, ss IEntity, intercepter ...ServerInterceptor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("RegisterService(%q)", sd.TypeName)
	if s.serve {
		fmt.Printf("grpc: EntityManager.RegisterService after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.entityDescs[sd.TypeName]; ok {
		fmt.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", sd.TypeName)
		return
	}
	sd.ServiceImpl = ss
	sd.interceptors = intercepter
	s.entityDescs[sd.TypeName] = sd
}

func (s *EntityManager) RegisterService(sd *ServiceDesc, ss IService, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			fmt.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss, intercepter...)
}

func (s *EntityManager) registerService(sd *ServiceDesc, ss IService, intercepter ...ServerInterceptor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("RegisterService(%q)", sd.TypeName)
	if s.serve {
		fmt.Printf("grpc: EntityManager.RegisterService after EntityManager.Serve for %q", sd.TypeName)
	}
	if _, ok := s.Services[sd.TypeName]; ok {
		fmt.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", sd.TypeName)
		return
	}
	sd.ServiceImpl = ss
	sd.interceptors = intercepter
	ss.setDesc(sd)
	ctx := allocContext()
	ctx.Service = ss
	ss.onInit(ctx)
	s.Services[sd.TypeName] = ss

	// start msg loop
	go ss.messageLoop()

	if ProxyRegister != nil {
		ProxyRegister.RegisterServiceToProxy(sd.TypeName)
	}
}
