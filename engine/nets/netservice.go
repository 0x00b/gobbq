package nets

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
)

type NetWorkName string

const (
	WebSocket NetWorkName = "websocket"
	TCP       NetWorkName = "tcp"
	TCP6      NetWorkName = "tcp6"
	KCP       NetWorkName = "kcp"
)

type NetName interface {
	// Name returns the name of the Transport implementation.
	// the result cannot change between calls.
	Name() NetWorkName
}

type ServiceName interface {
	// Name returns the name of the Transport implementation.
	// the result cannot change between calls.
	Name() string
}

// ServerTransport is the common interface for all gRPC server-side transport
// implementations.
//
// Methods may be called concurrently from multiple goroutines, but
// Write methods for a given Packet will be called serially.

type NetService interface {
	ServiceName

	ListenAndServe() error

	Close(chan struct{}) error
}

type service struct {
	opts *Options

	ctx    context.Context
	cancel func()

	// idleTimeout time.Duration
	// lastVisited time.Time

	closed atomic.Bool

	connMtx sync.Mutex
	conns   map[*Conn]struct{}

	websocket *WebSocketService
}

func NewNetService(opts ...Option) *service {
	svc := &service{
		opts:  &Options{},
		conns: make(map[*Conn]struct{}),
	}

	svc.ctx, svc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(svc.opts)
	}

	return svc
}

func (s *service) ListenAndServe() error {
	return s.listenAndServe(s.opts.network, s.opts.address, s.opts)
}

func (s *service) Close(closeChan chan struct{}) error {

	xlog.Infoln("server closing", s.opts.network, s.opts.address)

	s.cancel()

	defer s.closed.Store(true)

	if s.websocket != nil {
		err := s.websocket.Close(closeChan)
		if err != nil {
			return err
		}
	}

	s.closeAll()

	closeChan <- struct{}{}

	return nil
}

func (s *service) Name() string {
	return "service"
}

// ===== inner =====

func (s *service) listenAndServe(network NetWorkName, address string, opts *Options) error {

	if network == WebSocket {
		s.websocket = newWebSocketService(s)
		return s.websocket.ListenAndServe(network, address, opts)
	}

	var ln net.Listener
	var err error

	switch network {
	case KCP:
		ln, err = NewDefaultKCPListener(opts).Listen(network, address)
	case TCP, TCP6:
		ln, err = NewTCPListener(network, opts).Listen(network, address)
	default:
		panic(fmt.Sprintf("unkown network:%s", network))
	}

	if err != nil {
		return err
	}
	xlog.Infoln("listenAndServe from:", network, address)

	var once sync.Once
	closeListener := func() {
		if err := ln.Close(); err != nil {
			xlog.Error("listener close err:%s", err.Error())
		}
	}
	defer once.Do(closeListener)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			if err != nil {
				continue
			} else {
				return err
			}
		}

		xlog.Infof("Connection from: %s", conn.RemoteAddr())

		s.handleConn(conn, opts)
	}
}

func (s *service) handleConn(rawConn net.Conn, opts *Options) {

	if s.closed.Load() {
		return
	}

	if s.opts.TLSCertFile != "" && s.opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.opts.TLSCertFile, s.opts.TLSKeyFile)
		if err != nil {
			xlog.Traceln(err, "load RSA key & certificate failed")
			return
		}
		tlsConfig := &tls.Config{
			//MinVersion:       tls.VersionTLS12,
			//CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			Certificates: []tls.Certificate{cert},
			//CipherSuites: []uint16{
			//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			//	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			//},
			//PreferServerCipherSuites: true,
		}
		tlsConn := tls.Server(rawConn, tlsConfig)
		rawConn = net.Conn(tlsConn)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cn := newDefaultConn(ctx, rawConn, opts)
	cn.cancel = cancel

	// con err handler
	cn.registerConErrHandler(s)

	s.storeConn(cn)

	secure.GO(cn.Serve)

}

const MaxCloseWaitTime = 10

func (s *service) closeAll() {
	if s.closed.Load() {
		return
	}

	// close all conn
	closeWaitTime := s.opts.MaxCloseWaitTime
	if closeWaitTime < MaxCloseWaitTime {
		closeWaitTime = MaxCloseWaitTime
	}

	ctx, cancel := context.WithTimeout(context.Background(), closeWaitTime)
	defer cancel()

	var wg sync.WaitGroup

	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	for cn := range s.conns {

		wg.Add(1)

		cn := cn
		secure.GO(func() {
			defer wg.Done()

			c := make(chan struct{}, 1)
			secure.GO(func() {
				cn.Close(c)
			})

			select {
			case <-c:
			case <-ctx.Done():
			}
		})
	}

	wg.Wait()
}

func (s *service) storeConn(cn *Conn) {
	if s.closed.Load() {
		return
	}
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	s.conns[cn] = struct{}{}
}

func (s *service) unstoreConn(cn *Conn) {
	if s.closed.Load() {
		return
	}
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	delete(s.conns, cn)
}

// ConnCallback

func (s *service) HandleClose(cn *Conn) {
	if cn == nil {
		return
	}
	s.unstoreConn(cn)
}

func (s *service) HandleEOF(cn *Conn) {
	s.HandleClose(cn)
}

func (s *service) HandleTimeOut(cn *Conn) {
	s.HandleClose(cn)
}

func (s *service) HandleFail(cn *Conn) {
	s.HandleClose(cn)
}
