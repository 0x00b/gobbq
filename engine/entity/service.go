package entity

import (
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/xlog"
)

type IService interface {
	IBaseEntity

	ServiceDesc() *EntityDesc
	SetServiceDesc(desc *EntityDesc)

	serviceType()
}

type Service struct {
	baseEntity

	desc *EntityDesc
}

func (s *Service) ServiceDesc() *EntityDesc {
	return s.desc
}

func (s *Service) SetServiceDesc(desc *EntityDesc) {
	s.desc = desc
}

func (s *Service) getEntityMgr() *EntityManager {
	return s.desc.EntityMgr
}

func (e *Service) serviceType() {}

func (e *Service) onInit(c Context, cancel func(), id EntityID) {
	e.context = c
	e.cancel = cancel
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 10000)
	e.localCallChan = make(chan *localCall, 10000)
	e.callback = make(map[string]Callback, 10000)
	e.respChan = make(chan *codec.Packet, 10000)

	e.timer.Init()
	e.ticker = time.NewTicker(GAME_SERVICE_TICK_INTERVAL)

	e.OnInit()
}

func (e *Service) Run(ch chan bool) {
	xlog.Traceln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	defer func() {
		wg.Wait()
		// todo unregister service, and svcentity

	}()

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
			go func(ctx Context, release releaseCtx, pkt *codec.Packet) {
				defer release()
				defer wg.Done()

				err := e.handleCallMethod(ctx, pkt, e.ServiceDesc())
				if err != nil {
					xlog.Errorln(err)
				}
			}(ctx, release, pkt)

		case lc := <-e.localCallChan:
			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			go func(ctx Context, release releaseCtx, lc *localCall) {
				defer release()
				defer wg.Done()

				err := e.handleLocalCallMethod(ctx, lc, e.ServiceDesc())
				if err != nil {
					xlog.Errorln(err)
				}
			}(ctx, release, lc)

		case pkt := <-e.respChan:
			xlog.Tracef("handle: %s", pkt.String())

			wg.Add(1)

			// 异步
			ctx, release := e.context.Copy()
			go func(ctx Context, release releaseCtx, pkt *codec.Packet) {
				defer release()
				defer wg.Done()

				e.handleMethodRsp(ctx, pkt)
			}(ctx, release, pkt)

		case <-e.ticker.C:
			e.timer.Tick()
			e.context.Entity().OnTick()
		}
	}
}
