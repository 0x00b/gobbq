package entity

import (
	"sync"

	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
)

type IService interface {
	IEntity

	ServiceDesc() *EntityDesc

	setServiceDesc(desc *EntityDesc)
	serviceType()
}

type Service struct {
	Entity
}

func (s *Service) ServiceDesc() *EntityDesc {
	return s.desc
}

func SetServiceDesc(s IService, desc *EntityDesc) {
	s.setEntityDesc(desc)
}

func (s *Service) setServiceDesc(desc *EntityDesc) {
	s.desc = desc
}

func (e *Service) serviceType() {}

func (e *Service) onInit(c Context, cancel func(), id EntityID) {
	e.initOnce.Do(func() {
		e.callChan = make(chan *nets.Packet, 10000)
		e.localCallChan = make(chan *localCall, 10000)
		e.callback = make(map[string]Callback, 10000)
		e.respChan = make(chan *nets.Packet, 10000)

		e.defaultInit(c, cancel, id)
	})
}

func (e *Service) run(ch chan bool) {
	done := true
	e.runOnce.Do(func() {
		done = false
	})
	if done {
		ch <- true
		close(ch)
		return
	}

	xlog.Traceln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	tempch := make(chan bool)
	defer func() {
		wg.Wait()
		// todo unregister service, and svcentity

		e.onDestroy()

		close(ch)
		close(tempch)
	}()
	// response
	secure.GO(func() {
		for {
			select {
			case tempch <- true:

			case <-e.context.Done():
				xlog.Debugln("ctx done", e)
				return

			case pkt := <-e.respChan:
				// xlog.Tracef("handle: %s", pkt.String())

				wg.Add(1)

				// 异步
				ctx, release := e.context.Copy()
				npkt := pkt
				secure.GO(func() {
					defer wg.Done()
					defer release()
					e.handleMethodRsp(ctx, npkt)
				})
			}
		}
	})

	// 上面的for执行了，在继续下面的for
	<-tempch

	// async request, responese
	for {
		select {
		case ch <- true:

		case <-e.context.Done():
			xlog.Traceln("ctx done", e)
			return

		case pkt := <-e.callChan:
			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			npkt := pkt
			secure.GO(func() {
				defer wg.Done()
				defer release()

				e.handleCallMethod(ctx, npkt, e.ServiceDesc())
			})

		case lc := <-e.localCallChan:
			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			tlc := lc
			secure.GO(func() {
				defer wg.Done()
				defer release()

				e.handleLocalCallMethod(ctx, tlc, e.ServiceDesc())
			})

		case <-e.ticker.C:
			secure.DO(func() {
				e.timer.Tick()
				e.context.Entity().OnTick()
			})
		}
	}
}
