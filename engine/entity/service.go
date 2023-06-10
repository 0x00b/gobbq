package entity

import (
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

type IService interface {
	IBaseEntity

	serviceType()
}

type Service struct {
	baseEntity
}

func (e *Service) serviceType() {}

func (e *Service) onInit(c Context, cancel func(), id *bbq.EntityID) {
	e.context = c
	e.cancel = cancel
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 10000)
	e.localCallChan = make(chan *localCall, 10000)
	e.callback = make(map[string]Callback, 10000)
	e.respChan = make(chan *codec.Packet, 10000)

	e.timer.Init()
	e.ticker = time.Tick(GAME_SERVICE_TICK_INTERVAL)

	e.OnInit()
}

func (e *Service) Run() {
	xlog.Traceln("start message loop", e.EntityID())

	wg := sync.WaitGroup{}

	defer func() {
		wg.Wait()
		// todo unregister service, and svcentity

	}()

	// async request, responese
	for {
		select {
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

				err := e.handleCallMethod(ctx, pkt)
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

				err := e.handleLocalCallMethod(ctx, lc)
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

		case <-e.ticker:
			e.timer.Tick()
			e.context.Entity().OnTick()
		}
	}
}
