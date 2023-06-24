package main

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
)

//	替client实现 watch接口

// SysWatch
func (g *Gate) SysWatch(c entity.Context, req *entity.WatchRequest) (*entity.WatchResponse, error) {
	xlog.Infoln("SysWatch", req.EntityID)
	pkt := c.Packet()

	wc := g.watcher[entity.DstEntity(pkt)]
	if wc == nil {
		wc = make(map[entity.EntityID]bool)
		g.watcher[entity.DstEntity(pkt)] = wc
	}
	wc[entity.EntityID(req.EntityID)] = true

	return &entity.WatchResponse{}, nil
}

// SysUnwatch
func (g *Gate) SysUnwatch(c entity.Context, req *entity.WatchRequest) (*entity.WatchResponse, error) {
	xlog.Infoln("SysUnwatch", req.EntityID)

	pkt := c.Packet()
	wc := g.watcher[entity.DstEntity(pkt)]
	delete(wc, entity.EntityID(req.EntityID))

	return &entity.WatchResponse{}, nil
}
