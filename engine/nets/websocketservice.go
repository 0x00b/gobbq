package nets

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
	"golang.org/x/net/websocket"
)

// curl --include --no-buffer --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: example.com:80" --header "Origin: http://example.com:80" --header "Sec-WebSocketService-Key: SGVsbG8sIHdvcmxkIQ==" --header "Sec-WebSocketService-Version: 13" localhost:80

type WebSocketService struct {
	hs http.Server
}

func newWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (ws *WebSocketService) ListenAndServe(network NetWorkName, address string, opts *Options) error {
	if network != WebSocket {
		return errors.New("not websocket")
	}

	xlog.Infoln("websocket listenAndServe from:", network, address)

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
	ws.hs.Shutdown(context.Background())
	return nil
}

func (ws *WebSocketService) Name() NetWorkName {
	return WebSocket
}

func (ws *WebSocketService) handleConn(rawConn net.Conn, opts *Options) {

	xlog.Traceln("handleconn")

	cn := newDefaultConn(context.Background())

	cn.rwc = rawConn
	cn.packetReadWriter = codec.NewPacketReadWriter(rawConn)
	cn.PacketHandler = opts.PacketHandler
	cn.opts = opts
	if opts.ConnHandler != nil {
		cn.ConnHandler = opts.ConnHandler
	}

	cn.Serve()

	return
}
