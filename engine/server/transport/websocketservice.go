package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/engine/server"
	"golang.org/x/net/websocket"
)

// curl --include --no-buffer --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: example.com:80" --header "Origin: http://example.com:80" --header "Sec-WebSocketService-Key: SGVsbG8sIHdvcmxkIQ==" --header "Sec-WebSocketService-Version: 13" localhost:80

type WebSocketService struct {
	server http.Server
}

func newWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (ws *WebSocketService) ListenAndServe(network server.NetWorkName, address string, opts *server.ServerOptions) error {
	if network != server.WebSocket {
		return errors.New("not websocket")
	}

	ws.server.Addr = address

	h := websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		ws.handleConn(conn, opts)
	})

	ws.server.Handler = h
	// http.Handle("/ws", h)

	if opts.TLSKeyFile == "" && opts.TLSCertFile == "" {
		ws.server.ListenAndServe()
	} else {
		ws.server.ListenAndServeTLS(opts.CACertFile, opts.TLSKeyFile)
	}
	return nil
}

func (ws *WebSocketService) Close(chan struct{}) error {
	ws.server.Shutdown(context.Background())
	return nil
}

func (ws *WebSocketService) Name() server.NetWorkName {
	return server.WebSocket
}

func (ws *WebSocketService) handleConn(rawConn net.Conn, opts *server.ServerOptions) {

	fmt.Println("handleconn")

	conn := &conn{
		rwc:              rawConn,
		ctx:              context.Background(),
		packetReadWriter: codec.NewPacketReadWriter(context.Background(), rawConn),
		PacketHandler:    NewServerPacketHandler(context.Background(), rawConn, opts),
		opts:             opts,
	}
	conn.Serve()

	return
}
