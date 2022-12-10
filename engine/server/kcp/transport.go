package kcp

import (
	"errors"
	"fmt"
	"net"

	"github.com/0x00b/gobbq/engine/server"
	"github.com/xtaci/kcp-go"
)

// curl --include \
//      --no-buffer \
//      --header "Connection: Upgrade" \
//      --header "Upgrade: websocket" \
//      --header "Host: example.com:80" \
//      --header "Origin: http://example.com:80" \
//      --header "Sec-KCPTransport-Key: SGVsbG8sIHdvcmxkIQ==" \
//      --header "Sec-KCPTransport-Version: 13" \
//      localhost:80

type KCPTransport struct {
}

func (ks *KCPTransport) ListenAndServe(network server.NetWorkName, address string, ops server.ServerOptions) error {
	if network != server.KCP {
		return errors.New("not websocket")
	}

	kcpListener, err := kcp.ListenWithOptions(address, nil, 10, 3)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Listening on KCP: %s ...", address)

	for {
		conn, err := kcpListener.AcceptKCP()
		if err != nil {
			fmt.Println(err)
			return err
		}
		ks.handleConn(conn)
	}
}

func (ks *KCPTransport) Close(chan struct{}) error {
	return nil
}

func (ks *KCPTransport) Name() server.NetWorkName {
	return server.KCP
}

func (ks *KCPTransport) handleConn(conn net.Conn) {

	fmt.Println("handleconn")
	return
}
