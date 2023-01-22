package entity

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/proto"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[ServiceType]*ServiceDesc),
	entityDescs: make(map[ServiceType]*ServiceDesc),
	Entities:    make(map[EntityID]*ServiceDesc),
}

// EntityManager manage entity lifecycle
type EntityManager struct {
	mu    sync.Mutex // guards following
	serve bool

	Services    map[ServiceType]*ServiceDesc // service name -> service info
	entityDescs map[ServiceType]*ServiceDesc // entity name -> entity info
	Entities    map[EntityID]*ServiceDesc    // entity id -> entity impl
}

func NewEntity(typ ServiceType) *proto.Entity {
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
	entity.OnInit()

	eid := EntityID(entity.Entity().ID)

	fmt.Println("register entity id:", eid)

	newDesc := *desc
	newDesc.ServiceImpl = entity

	Manager.Entities[eid] = &newDesc

	// start message loop

	// send to poxy
	// proxy.RegisterEntity(tity.Entity())

	return nil
}

func (s *EntityManager) RegisterEntity(sd *ServiceDesc, ss IEntity) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			fmt.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerEntity(sd, ss)
}

func (s *EntityManager) registerEntity(sd *ServiceDesc, ss IEntity) {
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
	s.entityDescs[sd.TypeName] = sd
}

func (s *EntityManager) RegisterService(sd *ServiceDesc, ss IService) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			fmt.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss)
}

func (s *EntityManager) registerService(sd *ServiceDesc, ss IService) {
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
	s.Services[sd.TypeName] = sd
}
