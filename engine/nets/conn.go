package nets

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
)

func newDefaultConn(ctx context.Context, rawConn net.Conn, opts *Options) *Conn {
	tc := &Conn{
		opts:          opts,
		ctx:           ctx,
		rwc:           rawConn,
		PacketHandler: opts.PacketHandler,
	}

	tc.packetReadWriter = NewPacketReadWriter(tc)

	if opts.ConnErrHandler != nil {
		tc.connErrHandlers = append(tc.connErrHandlers, opts.ConnErrHandler)
	}

	return tc
}

type Conn struct {
	opts *Options

	ctx    context.Context
	cancel func()

	rwc              net.Conn
	packetReadWriter *PacketReadWriter

	idleTimeout time.Duration
	lastVisited time.Time

	PacketHandler   PacketHandler
	connErrHandlers []ConnErrHandler

	closeOnce sync.Once
	closed    bool
}

func (st *Conn) Name() string {
	return "server"
}

func (st *Conn) close() (e error) {
	st.closeOnce.Do(func() {
		if st.cancel != nil {
			st.cancel()
		}

		e = st.rwc.Close()
		if e != nil {
			xlog.Errorln(e)
		}
		st.closed = true
	})
	return e
}

func (st *Conn) Close(closeChan chan struct{}) (e error) {
	e = st.close()
	closeChan <- struct{}{}
	return e
}

func (st *Conn) Serve() {
	defer st.close()

	for {
		// xlog.Traceln("serve 1")
		// 检查上游是否关闭
		select {
		case <-st.ctx.Done():
			xlog.Trace("context done...")
			return
		default:
		}
		// xlog.Traceln("serve 2")

		// st.idleTimeout = 10 * time.Second
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

		// xlog.Traceln("serve 3", utils.GoID())
		pkt, release, err := st.packetReadWriter.ReadPacket()
		// xlog.Traceln("serve 4")
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

		// xlog.Traceln("serve 5")
		err = st.handle(pkt, release)
		// xlog.Traceln("serve 6")
		if err != nil {
			xlog.Errorln("handle failed", err)
		}
	}
}

func (st *Conn) handle(pkt *Packet, release ReleasePkt) error {
	defer release()

	// todo report

	return st.PacketHandler.HandlePacket(pkt)
}

func (st *Conn) handleEOF(err error) {
	xlog.Infoln("transport: conn  EOF ", err)

	for _, v := range st.connErrHandlers {
		v.HandleEOF(st)
	}
}

func (st *Conn) handleTimeOut(err error) {
	xlog.Infoln("transport: conn  Time out ", err)

	for _, v := range st.connErrHandlers {
		v.HandleTimeOut(st)
	}
}

func (st *Conn) handleFail(err error) {
	xlog.Errorln("transport: conn serve ReadFrame fail ", err)

	for _, v := range st.connErrHandlers {
		v.HandleFail(st)
	}
}

func (st *Conn) registerConErrHandler(ConnErrHandler ConnErrHandler) {
	if ConnErrHandler == nil {
		return
	}

	st.connErrHandlers = append(st.connErrHandlers, ConnErrHandler)
}

// AsyncWritePacket async writes a packet, this method will never block
func (st *Conn) SendPacket(p *Packet) (err error) {
	if st.closed {
		return errors.New("conn closing")
	}
	defer func() {
		if err != nil {
			xlog.Traceln("close conn...")
			st.close()
			return
		}
	}()

	defer func() {
		if e := recover(); e != nil {
			err = errors.New("conn closing")
		}
	}()
	timeout := time.Second * 5
	if timeout == 0 {
		return writeFull(st.rwc, p.Serialize())
	} else {
		ch := make(chan error)
		secure.GO(func() {
			defer func() {
				ch <- err
			}()
			err = writeFull(st.rwc, p.Serialize())
		})
		select {
		case err = <-ch:
			return err
		case <-time.After(timeout):
			return errors.New("conn ErrWriteBlocking")
		}
	}
}
