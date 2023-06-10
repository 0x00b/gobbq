package main

import (
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

// RegisterProxy
func (p *Proxy) RegisterProxy(c entity.Context, req *proxypb.RegisterProxyRequest) (*proxypb.RegisterProxyResponse, error) {

	p.proxyMap[req.ProxyID] = c.Packet().Src
	svcs := []string{}
	for n := range p.Server.EntityMgr.Services {
		svcs = append(svcs, string(n))
	}

	for n := range p.svcMaps {
		svcs = append(svcs, string(n))
	}

	return &proxypb.RegisterProxyResponse{SvcNames: svcs}, nil
}

// SyncService
func (p *Proxy) SyncService(c entity.Context, req *proxypb.SyncServiceRequest) (*proxypb.SyncServiceResponse, error) {

	p.RegisterProxyService(req.SvcName, c.Packet().Src)

	return &proxypb.SyncServiceResponse{}, nil
}

// Ping
func (p *Proxy) Ping(c entity.Context, req *proxypb.PingPong) (*proxypb.PingPong, error) {

	return &proxypb.PingPong{}, nil
}

// RegisterInst
func (p *Proxy) RegisterInst(c entity.Context, req *proxypb.RegisterInstRequest) (*proxypb.RegisterInstResponse, error) {
	xlog.Traceln("register inst", req.String())
	p.registerInst(req.GetInstID(), c.Packet().Src)
	xlog.Traceln("register inst done", req.String())
	return &proxypb.RegisterInstResponse{ProxyID: p.EntityID().ID}, nil
}

// RegisterEntity
func (p *Proxy) RegisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	xlog.Debugln("register service:", req.ServiceName)
	p.registerService(req.ServiceName, c.Packet().Src)

	for id, prw := range p.proxyMap {
		_ = prw
		entity.SetRemoteEntityManager(c, prw)
		_, err := proxypb.NewProxyEtyEntityClient(&bbq.EntityID{ID: id, ProxyID: id}).SyncService(c, &proxypb.SyncServiceRequest{SvcName: req.ServiceName})
		if err != nil {
			xlog.Errorln("sync svc", err)
			return nil, err
		}
	}

	return &proxypb.RegisterServiceResponse{}, nil
}

// RegisterEntity
func (p *Proxy) UnregisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	return &proxypb.RegisterServiceResponse{}, nil
}
