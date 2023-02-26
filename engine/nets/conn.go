package nets

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

func newDefaultConn(ctx context.Context, rawConn net.Conn, opts *Options) *conn {
	tc := &conn{
		opts:             opts,
		ctx:              ctx,
		rwc:              rawConn,
		packetReadWriter: codec.NewPacketReadWriter(rawConn),
		PacketHandler:    opts.PacketHandler,
	}

	if opts.ConnErrHandler != nil {
		tc.connErrHandlers = append(tc.connErrHandlers, opts.ConnErrHandler)
	}

	return tc
}

type conn struct {
	opts *Options

	ctx    context.Context
	cancel func()

	rwc              net.Conn
	packetReadWriter *codec.PacketReadWriter

	idleTimeout time.Duration
	lastVisited time.Time

	PacketHandler   PacketHandler
	connErrHandlers []ConnErrHandler
}

func (st *conn) Name() string {
	return "server"
}

func (st *conn) Close(closeChan chan struct{}) error {

	if st.cancel != nil {
		st.cancel()
	}

	err := st.rwc.Close()
	if err != nil {
		return err
	}

	closeChan <- struct{}{}

	return nil
}

func (st *conn) WritePacket(pkt *codec.Packet) error {
	return st.packetReadWriter.WritePacket(pkt)
}

func (st *conn) Serve() {
	for {
		xlog.Traceln("serve 1")
		// 检查上游是否关闭
		select {
		case <-st.ctx.Done():
			return
		default:
		}
		xlog.Traceln("serve 2")

		if st.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(st.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				st.lastVisited = now
				err := st.rwc.SetReadDeadline(now.Add(st.idleTimeout))
				if err != nil {
					xlog.Traceln("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		xlog.Traceln("serve 3")
		pkt, release, err := st.packetReadWriter.ReadPacket()
		xlog.Traceln("serve 4")
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				st.handleEOF(err)
				return
			}
			if e, ok := err.(net.Error); ok && e.Timeout() { // 客户端超过空闲时间没有发包，服务端主动超时关闭
				st.handleTimeOut(err)
				return
			}
			st.handleFail(err)
			return
		}

		err = st.handle(pkt, release)
		xlog.Traceln("serve 5")
		if err != nil {
			xlog.Errorln("handle failed", err)
		}
	}
}

func (st *conn) handle(pkt *codec.Packet, release codec.ReleasePkt) error {
	defer release()

	// todo report

	return st.PacketHandler.HandlePacket(pkt)
}

func (st *conn) handleEOF(err error) {
	xlog.Infoln("transport: conn  EOF ", err)

	for _, v := range st.connErrHandlers {
		v.HandleEOF(st)
	}
}

func (st *conn) handleTimeOut(err error) {
	xlog.Infoln("transport: conn  Time out ", err)

	for _, v := range st.connErrHandlers {
		v.HandleTimeOut(st)
	}
}

func (st *conn) handleFail(err error) {
	xlog.Errorln("transport: conn serve ReadFrame fail ", err)

	for _, v := range st.connErrHandlers {
		v.HandleFail(st)
	}
}

func (st *conn) registerConErrHandler(ConnErrHandler ConnErrHandler) {
	if ConnErrHandler == nil {
		return
	}

	st.connErrHandlers = append(st.connErrHandlers, ConnErrHandler)
}
