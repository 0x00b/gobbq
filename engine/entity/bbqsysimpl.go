package entity

import "github.com/0x00b/gobbq/xlog"

// entity 接口实现

var _ BbqSysEntity = &Entity{}

func (e *Entity) Watch(id EntityID) {
	client := NewBbqSysClient(id)
	client.SysWatch(e.Context(), &WatchRequest{EntityID: uint64(e.EntityID())})
}

func (e *Entity) Unwatch(id EntityID) {
	client := NewBbqSysClient(id)
	client.SysUnwatch(e.Context(), &WatchRequest{EntityID: uint64(e.EntityID())})
}

// Receive 默认实现
func (e *Entity) OnNotify(w NotifyInfo) {

	xlog.Infoln("default receive", w.EntityID)

}

// rpc接口实现

// SysWatch
func (e *Entity) SysWatch(c Context, req *WatchRequest) error {
	xlog.Infoln("SysWatch", req.EntityID)
	e.watchers[EntityID(req.EntityID)] = true
	return nil
}

// SysUnwatch
func (e *Entity) SysUnwatch(c Context, req *WatchRequest) error {
	xlog.Infoln("SysUnwatch", req.EntityID)
	delete(e.watchers, EntityID(req.EntityID))
	return nil
}

// SysNotify
func (e *Entity) SysNotify(c Context, req *WatchRequest) error {

	xlog.Infoln("SysNotify", req.EntityID)
	c.Entity().OnNotify(NotifyInfo{EntityID: EntityID(req.EntityID)})

	return nil
}
