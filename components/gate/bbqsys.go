package main

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbqsys"
)

type GateEntity struct {
	// 不能用这个，需要用service to entity，新实现方式，把service转变成entity， 保留service得EntityID
	entity.Entity

	gate *Gate
}

// SysWatch
func (g *GateEntity) SysWatch(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	return nil, nil
}

// SysUnwatch
func (g *GateEntity) SysUnwatch(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	return nil, nil
}

// SysNotify
func (g *GateEntity) SysNotify(c entity.Context, req *bbqsys.WatchRequest) (*bbqsys.WatchResponse, error) {
	return nil, nil
}
