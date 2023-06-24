package entity

import "github.com/0x00b/gobbq/xlog"

// entity 接口实现

func (e *Entity) Watch(id EntityID) {
	client := NewBbqSysEntityClient(id)
	client.SysWatch(e.Context(), &WatchRequest{EntityID: uint64(e.EntityID())})
}

func (e *Entity) Unwatch(id EntityID) {
	client := NewBbqSysEntityClient(id)
	client.SysUnwatch(e.Context(), &WatchRequest{EntityID: uint64(e.EntityID())})
}

// Receive 默认实现
func (e *Entity) Receive(w WatchNotify) {

	xlog.Infoln("default receive", w.EntityID)

}

// rpc接口实现

// SysWatch
func (e *Entity) SysWatch(c Context, req *WatchRequest) (*WatchResponse, error) {
	xlog.Infoln("SysWatch", req.EntityID)
	e.watchers[EntityID(req.EntityID)] = true
	return &WatchResponse{}, nil
}

// SysUnwatch
func (e *Entity) SysUnwatch(c Context, req *WatchRequest) (*WatchResponse, error) {
	xlog.Infoln("SysUnwatch", req.EntityID)
	delete(e.watchers, EntityID(req.EntityID))
	return &WatchResponse{}, nil
}

// SysNotify
func (e *Entity) SysNotify(c Context, req *WatchRequest) (*WatchResponse, error) {

	xlog.Infoln("SysNotify", req.EntityID)
	c.Entity().Receive(WatchNotify{EntityID: EntityID(req.EntityID)})

	return &WatchResponse{}, nil
}
