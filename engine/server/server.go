package server

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/bbqsync"
	"github.com/0x00b/gobbq/engine/entity"
)

// NewSever return gobbq server
func NewServer(opts ...ServerOption) *Server {
	svr := &Server{
		quit: bbqsync.NewEvent(),
		done: bbqsync.NewEvent(),
	}

	return svr
}

// Server is a gobbq server to serve RPC requests.
type Server struct {
	opts ServerOptions

	mu sync.Mutex // guards following
	// conns contains all active server transports. It is a map keyed on a
	// listener address with the value being the set of active transports
	// belonging to that listener.
	serve bool
	// services map[string]*serviceInfo // service name -> service info

	quit    *bbqsync.Event
	done    *bbqsync.Event
	serveWG sync.WaitGroup // counts active Serve goroutines for GracefulStop

	service Service
}

type ServerOptions struct {
	Network     string
	Address     string
	CACertFile  string // ca证书
	TLSCertFile string // server证书
	TLSKeyFile  string // server秘钥

	maxSendPacketSize int
	writeBufferSize   int
	readBufferSize    int
	numServerWorkers  uint32
	connectionTimeout time.Duration

	Entities map[entity.EntityType]*entity.EntityDesc
}

var ErrServerStopped = errors.New("gobbq: the server has been stopped")
var ErrNoServive = errors.New("gobbq: no register service")
var ErrServerUnknown = errors.New("gobbq: the network is unknown")

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type ServerOption interface {
	apply(*ServerOptions)
}

func (s *Server) ListenAndServe(network NetWorkName, address string, ops ...ServerOption) error {
	if s.service == nil {
		return ErrNoServive
	}
	err := s.service.ListenAndServe(network, address, s.opts)

	return err
}

// RegisterEntity RegisterService registers a service and its implementation to the gRPC
// server. It is called from the IDL generated code. This must be called before
// invoking Serve. If ss is non-nil (for legacy code), its type is checked to
// ensure it implements sd.HandlerType.
func (s *Server) RegisterEntity(sd *entity.EntityDesc, ss interface{}) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			fmt.Printf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *entity.EntityDesc, ss interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("RegisterService(%q)", sd.TypeName)
	if s.serve {
		fmt.Printf("grpc: Server.RegisterService after Server.Serve for %q", sd.TypeName)
	}
	if _, ok := s.opts.Entities[sd.TypeName]; ok {
		fmt.Printf("grpc: Server.RegisterService found duplicate service registration for %q", sd.TypeName)
		return
	}
	s.opts.Entities[sd.TypeName] = sd
}

func (s *Server) RegisterService(t Service) {
	s.service = t
}
