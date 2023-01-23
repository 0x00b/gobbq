package nets

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/0x00b/gobbq/engine/codec"
	"golang.org/x/net/websocket"
)

// curl --include --no-buffer --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: example.com:80" --header "Origin: http://example.com:80" --header "Sec-WebSocketService-Key: SGVsbG8sIHdvcmxkIQ==" --header "Sec-WebSocketService-Version: 13" localhost:80

type WebSocketService struct {
	hs  http.Server
	ctx context.Context
}

func newWebSocketService(ctx context.Context) *WebSocketService {
	return &WebSocketService{ctx: ctx}
}

func (ws *WebSocketService) ListenAndServe(network NetWorkName, address string, opts *Options) error {
	if network != WebSocket {
		return errors.New("not websocket")
	}

	ws.hs.Addr = address

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
	ws.hs.Shutdown(ws.ctx)
	return nil
}

func (ws *WebSocketService) Name() NetWorkName {
	return WebSocket
}

func (ws *WebSocketService) handleConn(rawConn net.Conn, opts *Options) {

	fmt.Println("handleconn")

	conn := &conn{
		rwc:              rawConn,
		ctx:              ws.ctx,
		packetReadWriter: codec.NewPacketReadWriter(ws.ctx, rawConn),
		PacketHandler:    opts.PacketHandler,
		opts:             opts,
	}
	conn.Serve()

	return
}
