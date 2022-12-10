package tcp

import (
	"errors"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/server"
)

type TCPTransport struct {
	network server.NetWorkName
}

func NewTCPTransport(net server.NetWorkName) *TCPTransport {
	return &TCPTransport{net}
}

func (ts *TCPTransport) ListenAndServe(net server.NetWorkName, address string, ops server.ServerOptions) error {
	if net != ts.network {
		return fmt.Errorf("not %s", string(ts.network))
	}
	ts.network = net
	ts.listenAndServe(net, address, ops)
	return nil
}

func (ts *TCPTransport) Close(chan struct{}) error {
	return nil
}

func (ts *TCPTransport) Name() server.NetWorkName {
	return ts.network
}

// ==== . inner ===

func (ts TCPTransport) listenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {
	if network != server.TCP {
		return errors.New("not websocket")
	}

	ln, err := net.Listen("tcp", address)
	fmt.Printf("Listening on TCP: %s ...", address)

	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err != nil {
				continue
			} else {
				return err
			}
		}

		fmt.Printf("Connection from: %s", conn.RemoteAddr())
		go ts.handleConn(conn)
	}

	// s.mu.Unlock()

	// var tempDelay time.Duration // how long to sleep on accept failure
	// for {
	// 	rawConn, err := t.Accept()
	// 	if err != nil {
	// 		if ne, ok := err.(interface {
	// 			Temporary() bool
	// 		}); ok && ne.Temporary() {
	// 			if tempDelay == 0 {
	// 				tempDelay = 5 * time.Millisecond
	// 			} else {
	// 				tempDelay *= 2
	// 			}
	// 			if max := 1 * time.Second; tempDelay > max {
	// 				tempDelay = max
	// 			}
	// 			s.mu.Lock()
	// 			fmt.Printf("Accept error: %v; retrying in %v\n", err, tempDelay)
	// 			s.mu.Unlock()
	// 			timer := time.NewTimer(tempDelay)
	// 			select {
	// 			case <-timer.C:
	// 			case <-s.quit.Done():
	// 				timer.Stop()
	// 				return nil
	// 			}
	// 			continue
	// 		}
	// 		s.mu.Lock()
	// 		fmt.Printf("done serving; Accept = %v\n", err)
	// 		s.mu.Unlock()

	// 		if s.quit.HasFired() {
	// 			return nil
	// 		}
	// 		return err
	// 	}
	// 	tempDelay = 0
	// 	// Start a new goroutine to deal with rawConn so we don't stall this Accept
	// 	// loop goroutine.
	// 	//
	// 	// Make sure we account for the goroutine so GracefulStop doesn't nil out
	// 	// s.conns before this conn can be added.
	// 	s.serveWG.Add(1)
	// 	go func() {
	// 		s.handleServerTransport(address, t.NewServerTransport(rawConn))
	// 		s.serveWG.Done()
	// 	}()
	// }
}

func (ts *TCPTransport) handleConn(conn net.Conn) {

	fmt.Println("handleconn")
	return
}
