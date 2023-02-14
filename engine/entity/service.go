package entity

import (
	"github.com/0x00b/gobbq/engine/codec"
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

func (e *Service) onInit(c Context, id EntityID) {
	e.context = c
	e.entityID = id
	e.callChan = make(chan *codec.Packet, 10000)
	e.callback = make(map[string]Callback, 10000)
	e.respChan = make(chan *codec.Packet, 10000)

	e.OnInit()
}

func (e *Service) Run() {
	xlog.Println("start message loop", e.EntityID())

	go func() {
		for !e.done {
			select {
			case <-e.context.Done():
				xlog.Println("ctx done", e)
			case pkt := <-e.respChan:
				xlog.Printf("handle: %s", pkt.String())

				// 异步
				ctx, release := e.context.Copy()
				go func(ctx Context, release releaseCtx, pkt *codec.Packet) {
					defer release()
					e.handleMethodRsp(ctx, pkt)
				}(ctx, release, pkt)
			}
		}
	}()

	for !e.done {
		select {
		case <-e.context.Done():
			xlog.Println("ctx done", e)

		case pkt := <-e.callChan:
			xlog.Printf("handle: %s", pkt.String())

			// 异步
			ctx, release := e.context.Copy()
			go func(ctx Context, release releaseCtx, pkt *codec.Packet) {
				defer release()
				e.handleCallMethod(ctx, pkt)
			}(ctx, release, pkt)
		}
	}
	xlog.Println("stop message loop", e.EntityID())
	// todo unregister entity
}
