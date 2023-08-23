package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/0x00b/gobbq/components/proxy/proxypb"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/engine/nets"
	"github.com/0x00b/gobbq/xlog"
	capi "github.com/hashicorp/consul/api"
)

func (p *Proxy) RegisterToConsul() {

	ip := CFG.IP
	port := CFG.Port

	if CFG.RunEnv == "k8s" {
		// todo 环境变量获取
		ip = os.Getenv("SERVICE_<SERVICE-NAME>_SERVICE_HOST")
		port = os.Getenv("SERVICE_<SERVICE-NAME>_SERVICE_PORT")
	}

	tmpPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	serviceName := "gobbqproxy"

	pid := p.EntityID().ProxyID()
	check := &capi.AgentServiceCheck{
		TCP:      fmt.Sprintf("%s:%d", ip, tmpPort), // 这里一定是外部可以访问的地址
		Timeout:  "10s",                             // 超时时间
		Interval: "10s",                             // 运行检查的频率
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
		DeregisterCriticalServiceAfter: "1m",
	}
	srv := &capi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%d_%s-%d", serviceName, pid, ip, tmpPort), // 服务唯一ID
		Name:    serviceName,                                               // 服务名称
		Tags:    []string{"gobbqproxy"},                                    // 为服务打标签
		Meta:    map[string]string{"proxy_id": fmt.Sprint(pid)},
		Address: ip,
		Port:    tmpPort,
		Check:   check,
	}

	err = p.consul.Agent().ServiceRegister(srv)
	if err != nil {
		panic(err)
	}

	xlog.Infoln("register myself:", ip, port)
}

func (p *Proxy) ConnOtherProxys(ops ...nets.Option) {

	// svcs, err := p.consul.Agent().ServicesWithFilter("Service==`gobbqproxy`")
	_, svcs, err := p.consul.Agent().AgentHealthServiceByName("gobbqproxy")
	if err != nil {
		panic(err)
	}
	for _, svc := range svcs {

		p.ConnOtherProxy(&svc, ops...)

		xlog.Infoln("ConnOtherProxy:", svc)
	}
}

func (p *Proxy) ConnOtherProxy(pcfg *capi.AgentServiceChecksInfo, ops ...nets.Option) {

	svc := pcfg.Service

	prxy, err := nets.Connect(nets.TCP, svc.Address, fmt.Sprint(svc.Port), ops...)

	if err != nil {
		panic(err)
	}

	c, release := p.Context().Copy()
	// defer release()
	_ = release

	entity.SetProxy(c, prxy)

	id, ok := svc.Meta["proxy_id"]
	if !ok {
		panic(err)
	}
	eid, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	other := proxypb.NewProxyEtyClient(entity.FixedEntityID(entity.ProxyID(eid), 0, 0))

	rsp, err := other.RegisterProxy(c,
		&proxypb.RegisterProxyRequest{
			ProxyID: uint32(p.EntityID().ProxyID()),
		})
	if err != nil {
		panic(err)
	}
	// xlog.Println("RegisterProxy done...")

	// p.proxyMtx.Lock()
	// defer p.proxyMtx.Unlock()
	p.proxyMap[entity.ProxyID(eid)] = prxy.GetConn()

	for _, v := range rsp.ServiceNames {
		p.RegisterProxyService(v, prxy.GetConn())
	}

	p.AddTimer(10*time.Second, func() bool {
		other.Ping(c, &proxypb.PingPong{})
		return true
	})
}
