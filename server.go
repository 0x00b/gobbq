package bbq

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
)

// NewSever return gobbq server
func NewServer() *Server {

	svr := &Server{
		EntityMgr: entity.NewEntityManager(),
	}

	return svr
}

// Server is a gobbq server to serve RPC requests.
type Server struct {
	netsvc []nets.NetService

	MaxCloseWaitTime time.Duration // max waiting time when closing server

	mux sync.Mutex // guards onShutdownHooks
	// onShutdownHooks are hook functions that would be executed when server is
	// shutting down (before closing all services of the server).
	onShutdownHooks []func()

	signalCh  chan os.Signal
	closeCh   chan struct{}
	closeOnce sync.Once

	EntityMgr *entity.EntityManager
}

var ErrServerStopped = errors.New("gobbq: the server has been stopped")
var ErrNoServive = errors.New("gobbq: no register service")
var ErrServerUnknown = errors.New("gobbq: the network is unknown")

func (s *Server) ListenAndServe() error {
	if s.netsvc == nil {
		return ErrNoServive
	}
	for _, ns := range s.netsvc {
		err := ns.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// Close implements Service interface, notifying all services of server shutdown.
// Would wait no more than 10s.
func (s *Server) Close(ch chan struct{}) error {
	if s.closeCh != nil {
		close(s.closeCh)
	}

	s.tryClose()

	if ch != nil {
		ch <- struct{}{}
	}
	return nil
}

// MaxCloseWaitTime is the max waiting time for closing services.
const MaxCloseWaitTime = 10 * time.Second

func (s *Server) tryClose() {
	fn := func() {
		// execute shutdown hook functions before closing services.
		s.mux.Lock()
		for _, f := range s.onShutdownHooks {
			f()
		}
		s.mux.Unlock()

		// close all Services
		closeWaitTime := s.MaxCloseWaitTime
		if closeWaitTime < MaxCloseWaitTime {
			closeWaitTime = MaxCloseWaitTime
		}
		ctx, cancel := context.WithTimeout(context.Background(), closeWaitTime)
		defer cancel()

		var wg sync.WaitGroup

		wg.Add(1)
		// close entity manager
		go func() {
			defer wg.Done()

			c := make(chan struct{}, 1)
			s.EntityMgr.Close(c)

			select {
			case <-c:
			case <-ctx.Done():
			}
		}()

		// close conn
		for _, service := range s.netsvc {
			wg.Add(1)
			go func(srv nets.NetService) {
				defer wg.Done()

				c := make(chan struct{}, 1)
				go srv.Close(c)

				select {
				case <-c:
				case <-ctx.Done():
				}
			}(service)
		}
		wg.Wait()
	}
	s.closeOnce.Do(fn)
}

// RegisterOnShutdown registers a hook function that would be executed when server is shutting down.
func (s *Server) RegisterOnShutdown(fn func()) {
	if fn == nil {
		return
	}
	s.mux.Lock()
	s.onShutdownHooks = append(s.onShutdownHooks, fn)
	s.mux.Unlock()
}

func (s *Server) RegisterNetService(ns nets.NetService) {
	s.netsvc = append(s.netsvc, ns)
}
