package entity

import (
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
)

type IService interface {
	IEntity

	ServiceDesc() *EntityDesc
	SetServiceDesc(desc *EntityDesc)

	serviceType()
}

type Service struct {
	Entity
}

func (s *Service) ServiceDesc() *EntityDesc {
	return s.desc
}

func (s *Service) SetServiceDesc(desc *EntityDesc) {
	s.desc = desc
}

func (e *Service) serviceType() {}

func (e *Service) onInit(c Context, cancel func(), id EntityID) {
	e.initOnce.Do(func() {
		e.context = c
		e.cancel = cancel
		e.entityID = id
		e.callChan = make(chan *nets.Packet, 10000)
		e.localCallChan = make(chan *localCall, 10000)
		e.callback = make(map[string]Callback, 10000)
		e.respChan = make(chan *nets.Packet, 10000)
		e.watchers = make(map[EntityID]bool)

		e.timer.Init()
		e.ticker = time.NewTicker(GAME_SERVICE_TICK_INTERVAL)

		e.OnInit()
	})
}

func (e *Service) Run() {
	e.runOnce.Do(func() {
		ch := make(chan bool)
		secure.GO(func() {
			e.run(ch)
		})
		<-ch
	})
}

func (e *Service) run(ch chan bool) {
	xlog.Traceln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	defer func() {
		wg.Wait()
		// todo unregister service, and svcentity

		e.onDestroy()

	}()

	tempch := make(chan bool)
	// response
	secure.GO(func() {
		for {
			select {
			case tempch <- true:

			case <-e.context.Done():
				xlog.Debugln("ctx done", e)
				return

			case pkt := <-e.respChan:
				xlog.Tracef("handle: %s", pkt.String())

				wg.Add(1)

				// 异步
				ctx, release := e.context.Copy()
				npkt := pkt
				secure.GO(func() {
					defer release()
					defer wg.Done()
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
			xlog.Tracef("handle: %s", pkt.String())

			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			npkt := pkt
			secure.GO(func() {
				defer release()
				defer wg.Done()

				err := e.handleCallMethod(ctx, npkt, e.ServiceDesc())
				if err != nil {
					xlog.Errorln(err)
				}
			})

		case lc := <-e.localCallChan:
			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			tlc := lc
			secure.GO(func() {
				defer release()
				defer wg.Done()

				err := e.handleLocalCallMethod(ctx, tlc, e.ServiceDesc())
				if err != nil {
					xlog.Errorln(err)
				}
			})

		case <-e.ticker.C:
			e.timer.Tick()
			e.context.Entity().OnTick()
		}
	}
}
