package entity

import (
	"reflect"
	"sync"

	"github.com/0x00b/gobbq/xlog"
)

var Manager EntityManager = EntityManager{
	Services:    make(map[TypeName]IService),
	entityDescs: make(map[TypeName]*EntityDesc),
	Entities:    make(map[string]IBaseEntity),
}
var ProxyRegister RegisterProxy

type RegisterProxy interface {
	RegisterEntityToProxy(eid EntityID) error
	RegisterServiceToProxy(svcName TypeName) error

	// UnregisterEntityToProxy(eid EntityID) error
	// UnregisterServiceToProxy(svcName TypeName) error
}

// EntityManager manage entity lifecycle
type EntityManager struct {
	mu    sync.RWMutex // guards following
	serve bool

	Services    map[TypeName]IService    // service name -> service info
	entityDescs map[TypeName]*EntityDesc // entity name -> entity info
	Entities    map[string]IBaseEntity   // entity id -> entity impl

}

func RegisterEntity(c Context, id *EntityID, entity IBaseEntity) error {
	ctx := allocContext()
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

func NewEntity(c Context, id *EntityID) error {
	desc, ok := Manager.entityDescs[id.Type]
	if !ok {
		xlog.Printf("grpc: EntityManager.RegisterService found duplicate service registration for %q", id.Type)
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

	xlog.Println("register entity id:", *id)

	newDesc := *desc
	newDesc.EntityImpl = entity
	entity.setDesc(&newDesc)

	RegisterEntity(c, id, entity)

	// start message loop
	go entity.Run()

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

func (s *EntityManager) RegisterService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			xlog.Printf("grpc: EntityManager.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.registerService(sd, ss, intercepter...)
}

func (s *EntityManager) registerService(sd *EntityDesc, ss IService, intercepter ...ServerInterceptor) {
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

	s.registerServiceEntity(sd, ss)

	xlog.Printf("grpc: EntityManager.RegisterService eid:%s", ss.EntityID())

	// start msg loop
	go ss.Run()

	if ProxyRegister != nil {
		ProxyRegister.RegisterServiceToProxy(sd.TypeName)
		ProxyRegister.RegisterEntityToProxy(ss.EntityID())
	}
}

func (s *EntityManager) registerServiceEntity(sd *EntityDesc, entity IService) error {

	RegisterEntity(nil, NewEntityID.NewEntityID(sd.TypeName), entity)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	s.Services[sd.TypeName] = entity

	return nil
}
