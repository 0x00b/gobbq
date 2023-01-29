package entity

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/components/proxy/ex"
	"github.com/0x00b/gobbq/proto/bbq"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[TypeName]*ServiceDesc),
	entityDescs: make(map[TypeName]*ServiceDesc),
	Entities:    make(map[EntityID]*ServiceDesc),
}

// EntityManager manage entity lifecycle
type EntityManager struct {
	mu    sync.Mutex // guards following
	serve bool

	Services    map[TypeName]*ServiceDesc // service name -> service info
	entityDescs map[TypeName]*ServiceDesc // entity name -> entity info
	Entities    map[EntityID]*ServiceDesc // entity id -> entity impl
}

func NewEntity(id EntityID, typ TypeName) *bbq.EntityID {
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
	entity.SetEntityID(id)
	// init
	entity.OnInit()

	eid := EntityID(entity.EntityID())

	fmt.Println("register entity id:", eid)

	newDesc := *desc
	newDesc.ServiceImpl = entity

	Manager.Entities[eid] = &newDesc

	// start message loop

	// send to poxy
	ex.RegisterEntity(string(id))

	return nil
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
	s.Services[sd.TypeName] = sd
}
