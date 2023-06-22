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

	if opts.ConnCallback != nil {
		tc.ConnCallbacks = append(tc.ConnCallbacks, opts.ConnCallback)
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

	PacketHandler PacketHandler
	ConnCallbacks []ConnCallback

	closeOnce sync.Once
	closed    bool
}

func (cn *Conn) Name() string {
	return "server"
}

func (cn *Conn) close() (e error) {
	cn.closeOnce.Do(func() {
		if cn.cancel != nil {
			cn.cancel()
		}

		e = cn.rwc.Close()
		if e != nil {
			xlog.Errorln(e)
		}
		cn.closed = true
	})
	return e
}

func (cn *Conn) Close(closeChan chan struct{}) (e error) {
	e = cn.close()
	closeChan <- struct{}{}
	return e
}

func (cn *Conn) Serve() {
	defer cn.close()

	for {
		// xlog.Traceln("serve 1")
		// 检查上游是否关闭
		select {
		case <-cn.ctx.Done():
			xlog.Trace("context done...")
			return
		default:
		}
		// xlog.Traceln("serve 2")

		// cn.idleTimeout = 10 * time.Second
		if cn.idleTimeout > 0 {
			now := time.Now()
			if now.Sub(cn.lastVisited) > 5*time.Second { // SetReadDeadline性能损耗较严重，每5s才更新一次timeout
				cn.lastVisited = now
				err := cn.rwc.SetReadDeadline(now.Add(cn.idleTimeout))
				if err != nil {
					xlog.Traceln("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		// xlog.Traceln("serve 3", utils.GoID())
		pkt, release, err := cn.packetReadWriter.ReadPacket()
		// xlog.Traceln("serve 4")
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				cn.handleEOF(err)
				return
			}
			if e, ok := err.(net.Error); ok && e.Timeout() { // 客户端超过空闲时间没有发包，服务端主动超时关闭
				cn.handleTimeOut(err)
				return
			}
			cn.handleFail(err)
			return
		}

		// xlog.Traceln("serve 5")
		err = cn.handle(pkt, release)
		// xlog.Traceln("serve 6")
		if err != nil {
			xlog.Errorln("handle failed", err)
		}
	}
}

func (cn *Conn) handle(pkt *Packet, release ReleasePkt) error {
	defer release()

	// todo report

	return cn.PacketHandler.HandlePacket(pkt)
}

func (cn *Conn) handleEOF(err error) {
	xlog.Infoln("transport: conn  EOF ", err)

	for _, v := range cn.ConnCallbacks {
		v.HandleEOF(cn)
	}
}

func (cn *Conn) handleTimeOut(err error) {
	xlog.Infoln("transport: conn  Time out ", err)

	for _, v := range cn.ConnCallbacks {
		v.HandleTimeOut(cn)
	}
}

func (cn *Conn) handleFail(err error) {
	xlog.Errorln("transport: conn serve ReadFrame fail ", err)

	for _, v := range cn.ConnCallbacks {
		v.HandleFail(cn)
	}
}

func (cn *Conn) registerConErrHandler(ConnCallback ConnCallback) {
	if ConnCallback == nil {
		return
	}

	cn.ConnCallbacks = append(cn.ConnCallbacks, ConnCallback)
}

// AsyncWritePacket async writes a packet, this method will never block
func (cn *Conn) SendPacket(p *Packet) (err error) {
	if cn.closed {
		return errors.New("conn closing")
	}
	defer func() {
		if err != nil {
			xlog.Traceln("close conn...")
			cn.close()
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
		return writeFull(cn.rwc, p.Serialize())
	} else {
		ch := make(chan error)
		secure.GO(func() {
			defer func() {
				ch <- err
			}()
			err = writeFull(cn.rwc, p.Serialize())
		})
		select {
		case err = <-ch:
			return err
		case <-time.After(timeout):
			return errors.New("conn ErrWriteBlocking")
		}
	}
}
