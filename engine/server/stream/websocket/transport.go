package websocket

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/0x00b/gobbq/engine/server"
	"golang.org/x/net/websocket"
)

// curl --include \
//      --no-buffer \
//      --header "Connection: Upgrade" \
//      --header "Upgrade: websocket" \
//      --header "Host: example.com:80" \
//      --header "Origin: http://example.com:80" \
//      --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
//      --header "Sec-WebSocket-Version: 13" \
//      localhost:80

type WebSocket struct {
	server http.Server
}

func (ws *WebSocket) ListenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {
	if network != server.WebSocket {
		return errors.New("not websocket")
	}

	ws.server.Addr = address

	h := websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		ws.handleConn(conn)
	})

	ws.server.Handler = h
	// http.Handle("/ws", h)

	if ops.TLSKeyFile == "" && ops.TLSCertFile == "" {
		ws.server.ListenAndServe()
	} else {
		ws.server.ListenAndServeTLS(ops.CACertFile, ops.TLSKeyFile)
	}
	return nil
}

func (ws *WebSocket) Close(chan struct{}) error {
	ws.server.Shutdown(context.Background())
	return nil
}

func (ws *WebSocket) Name() server.NetWorkName {
	return server.WebSocket
}

func (ws *WebSocket) handleConn(conn net.Conn) {

	fmt.Println("handleconn")
	return
}
