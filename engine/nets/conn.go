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
		pktRecvChan:   make(chan *Packet, 500),
		pktSendChan:   make(chan *Packet, 500),
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

	pktRecvChan chan *Packet
	pktSendChan chan *Packet
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

		close(cn.pktRecvChan)
		close(cn.pktSendChan)

		cn.closed = true
	})
	return e
}

func (cn *Conn) Close(closeChan chan struct{}) (e error) {
	e = cn.close()
	if closeChan != nil {
		closeChan <- struct{}{}
	}
	return e
}

func (cn *Conn) Serve() {
	defer cn.close()

	secure.GO(cn.handleLoop)
	secure.GO(cn.writeLoop)

	// recv loop
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

		cn.idleTimeout = 30 * time.Second
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
		pkt, err := cn.packetReadWriter.ReadPacket()
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

		cn.pktRecvChan <- pkt
	}
}

func (cn *Conn) handle(pkt *Packet) error {
	defer pkt.Release()

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
	xlog.Traceln("start send:", p.String())

	if cn.closed {
		return errors.New("conn closed")
	}
	defer func() {
		if err != nil {
			xlog.Traceln("close conn...")
			cn.close()
			return
		}
	}()

	// 需要改写,可配置
	timeout := time.Second * 5

	if timeout == 0 {
		select {
		case <-cn.ctx.Done():
			xlog.Trace("context done...")
			return

		case cn.pktSendChan <- p:
			p.Retain()
			return nil

		default:
			return errors.New("conn blocking")
		}

	} else {
		select {
		case <-cn.ctx.Done():
			xlog.Trace("context done...")
			return

		case cn.pktSendChan <- p:
			p.Retain()
			return nil

		case <-time.After(timeout):
			return errors.New("conn timeount blocking")
		}
	}
}

func (cn *Conn) writeLoop() {
	defer func() {
		cn.close()
	}()

	for {
		select {
		case <-cn.ctx.Done():
			xlog.Trace("context done...")
			return

		case p := <-cn.pktSendChan:
			func() {
				defer p.Release()
				if cn.closed {
					// 需要改写, 可能chan中还有pkt,需要pkt.Release
					return
				}
				cn.rwc.SetWriteDeadline(time.Now().Add(5 * time.Second))
				err := writeFull(cn.rwc, p.Serialize())
				if err != nil {
					xlog.Errorln(err)
					return
				}
			}()
		}
	}
}

func (cn *Conn) handleLoop() {
	defer func() {
		cn.close()
	}()

	for {
		select {
		case <-cn.ctx.Done():
			xlog.Trace("context done...")
			return

		case pkt := <-cn.pktRecvChan:
			if cn.closed {
				// 需要改写, 可能chan中还有pkt,需要pkt.Release
				return
			}

			err := cn.handle(pkt)
			if err != nil {
				xlog.Errorln("handle failed", err)
			}
		}
	}
}
