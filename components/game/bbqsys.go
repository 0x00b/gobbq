package game

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbqsys"
)

// SysWatch
func (g *Game) SysWatch(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	// todo 找到entity
	return nil, nil
}

// SysUnwatch
func (g *Game) SysUnwatch(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	return nil, nil
}

// SysNotify
func (g *Game) SysNotify(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	return nil, nil
}
