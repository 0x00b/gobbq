package main

import (
	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/xlog"
)

type ProxyService struct {
	entity.Service
}

type ProxyEntity struct {
	entity.Entity
}

func (ps *ProxyService) OnInit() {
	xlog.Println("on init ProxyService")
}

// RegisterProxy
func (ps *ProxyEntity) RegisterProxy(c entity.Context, req *proxypb.RegisterProxyRequest) (*proxypb.RegisterProxyResponse, error) {

	proxyInst.proxyMap[req.ProxyID] = c.Packet().Src
	svcs := []string{}
	for n := range entity.Manager.Services {
		svcs = append(svcs, string(n))
	}

	for n := range proxyInst.svcMaps {
		svcs = append(svcs, string(n))
	}

	return &proxypb.RegisterProxyResponse{SvcNames: svcs}, nil
}

// SyncService
func (ps *ProxyEntity) SyncService(c entity.Context, req *proxypb.SyncServiceRequest) (*proxypb.SyncServiceResponse, error) {

	proxyInst.RegisterProxyService(req.SvcName, c.Packet().Src)

	return &proxypb.SyncServiceResponse{}, nil
}

// Ping
func (ps *ProxyEntity) Ping(c entity.Context, req *proxypb.PingPong) (*proxypb.PingPong, error) {

	return &proxypb.PingPong{}, nil
}

// RegisterInst
func (ps *ProxyService) RegisterInst(c entity.Context, req *proxypb.RegisterInstRequest) (*proxypb.RegisterInstResponse, error) {

	return &proxypb.RegisterInstResponse{ProxyID: proxyInst.EntityID().ID}, nil
}

// RegisterEntity
func (ps *ProxyService) RegisterEntity(c entity.Context, req *proxypb.RegisterEntityRequest) (*proxypb.RegisterEntityResponse, error) {

	proxyInst.RegisterEntity(req.EntityID, c.Packet().Src)

	return &proxypb.RegisterEntityResponse{}, nil
}

// RegisterEntity
func (ps *ProxyService) RegisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	xlog.Println("register service:", req.ServiceName)
	proxyInst.RegisterService(req.ServiceName, c.Packet().Src)

	for id, p := range proxyInst.proxyMap {
		_, err := proxypb.NewProxyEtyEntityClient(p, &bbq.EntityID{ID: id, ProxyID: id}).SyncService(c, &proxypb.SyncServiceRequest{SvcName: req.ServiceName})
		if err != nil {
			xlog.Errorln("sync svc", err)
			return nil, err
		}
	}

	return &proxypb.RegisterServiceResponse{}, nil
}

// UnregisterEntity
func (ps *ProxyService) UnregisterEntity(c entity.Context, req *proxypb.RegisterEntityRequest) (*proxypb.RegisterEntityResponse, error) {

	return &proxypb.RegisterEntityResponse{}, nil
}

// RegisterEntity
func (ps *ProxyService) UnregisterService(c entity.Context, req *proxypb.RegisterServiceRequest) (*proxypb.RegisterServiceResponse, error) {

	return &proxypb.RegisterServiceResponse{}, nil
}

// Ping
func (ps *ProxyService) Ping(c entity.Context, req *proxypb.PingPong) (*proxypb.PingPong, error) {

	return &proxypb.PingPong{}, nil
}
