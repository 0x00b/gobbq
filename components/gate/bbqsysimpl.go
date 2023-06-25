package main

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
)

//	替client实现 watch接口

// SysWatch
func (g *Gate) SysWatch(c entity.Context, req *entity.WatchRequest) (*entity.WatchResponse, error) {
	pkt := c.Packet()
	dst := entity.DstEntity(pkt)

	xlog.Infoln("SysWatch", req.EntityID)

	if g.IsMyEntity(dst) {
		return g.Service.SysWatch(c, req)
	}

	wc := g.watcher[dst]
	if wc == nil {
		wc = make(map[entity.EntityID]bool)
		g.watcher[dst] = wc
	}
	wc[entity.EntityID(req.EntityID)] = true

	return &entity.WatchResponse{}, nil
}

// SysUnwatch
func (g *Gate) SysUnwatch(c entity.Context, req *entity.WatchRequest) (*entity.WatchResponse, error) {
	xlog.Infoln("SysUnwatch", req.EntityID)

	pkt := c.Packet()
	dst := entity.DstEntity(pkt)

	if g.IsMyEntity(dst) {
		return g.Service.SysUnwatch(c, req)
	}

	wc := g.watcher[dst]
	delete(wc, entity.EntityID(req.EntityID))

	return &entity.WatchResponse{}, nil
}
