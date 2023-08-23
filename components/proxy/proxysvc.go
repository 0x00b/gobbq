package main

import (
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/xlog"
)

type ProxySvc struct {
	entity.Service
	proxy *Proxy
}

// RegisterInst
func (p *ProxySvc) RegisterInst(c entity.Context, req *proxypb.RegisterInstRequest) (*proxypb.RegisterInstResponse, error) {
	xlog.Traceln("register inst", req.String())

	instID := p.proxy.NewInstID()
	p.proxy.registerInst(instID, c.Packet().Src)

	xlog.Traceln("register inst done", req.String())

	return &proxypb.RegisterInstResponse{
		ProxyID: uint32(p.EntityID().ProxyID()),
		InstID:  uint32(instID),
	}, nil
}

// Ping
func (p *ProxySvc) Ping(c entity.Context, req *proxypb.PingPong) (*proxypb.PingPong, error) {

	return &proxypb.PingPong{}, nil
}
