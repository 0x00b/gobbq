package server

import (
	"errors"
	"sync"

	"github.com/0x00b/gobbq/engine/bbqsync"
)

// NewSever return gobbq server
func NewServer(opts ...ServerOption) *Server {

	svr := &Server{
		quit: bbqsync.NewEvent(),
		done: bbqsync.NewEvent(),
		opts: &ServerOptions{},
	}

	for _, opt := range opts {
		opt.apply(svr.opts)
	}

	return svr
}

// Server is a gobbq server to serve RPC requests.
type Server struct {
	opts *ServerOptions

	quit    *bbqsync.Event
	done    *bbqsync.Event
	serveWG sync.WaitGroup // counts active Serve goroutines for GracefulStop

	netService NetService
}

var ErrServerStopped = errors.New("gobbq: the server has been stopped")
var ErrNoServive = errors.New("gobbq: no register service")
var ErrServerUnknown = errors.New("gobbq: the network is unknown")

func (s *Server) ListenAndServe(network NetWorkName, address string, ops ...ServerOption) error {
	if s.netService == nil {
		return ErrNoServive
	}
	err := s.netService.ListenAndServe(network, address, s.opts)

	return err
}

func (s *Server) RegisterNetService(t NetService) {
	s.netService = t
}
