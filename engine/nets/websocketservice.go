package nets

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/0x00b/gobbq/xlog"
	"golang.org/x/net/websocket"
)

// curl --include --no-buffer --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: example.com:80" --header "Origin: http://example.com:80" --header "Sec-WebSocketService-Key: SGVsbG8sIHdvcmxkIQ==" --header "Sec-WebSocketService-Version: 13" localhost:80

type WebSocketService struct {
	hs  http.Server
	svc *service
}

func newWebSocketService(svc *service) *WebSocketService {
	return &WebSocketService{
		svc: svc,
	}
}

func (ws *WebSocketService) ListenAndServe(network NetWorkName, address string, opts *Options) error {
	if network != WebSocket {
		return errors.New("not websocket")
	}

	xlog.Infoln("websocket listenAndServe from:", network, address)

	ws.hs.Addr = address
	ws.hs.SetKeepAlivesEnabled(opts.NetKeepAlive)
	ws.hs.RegisterOnShutdown(func() {

	})

	h := websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		ws.handleConn(conn, opts)
	})

	ws.hs.Handler = h
	// http.Handle("/ws", h)

	if opts.TLSKeyFile == "" && opts.TLSCertFile == "" {
		return ws.hs.ListenAndServe()
	}

	return ws.hs.ListenAndServeTLS(opts.CACertFile, opts.TLSKeyFile)
}

func (ws *WebSocketService) Close(chan struct{}) error {

	err := ws.hs.Shutdown(context.Background())
	return err
}

func (ws *WebSocketService) Name() NetWorkName {
	return WebSocket
}

func (ws *WebSocketService) handleConn(rawConn net.Conn, opts *Options) {
	if ws.svc.closed.Load() {
		xlog.Infoln("closed", ws.hs.Addr)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	cn := newDefaultConn(ctx, rawConn, opts)
	cn.cancel = cancel

	// con err handler
	cn.registerConErrHandler(ws.svc)

	ws.svc.storeConn(cn)

	cn.Serve()
}
