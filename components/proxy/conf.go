package main

import (
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/db/mongo"
	"github.com/0x00b/gobbq/xlog"
)

type Config struct {
	Env                    string `yaml:"env"`                      //  env: product/test/dev
	RunEnv                 string `yaml:"run_env"`                  //  run_env: host,k8s
	LogFile                string `yaml:"log_file"`                 //	log_file: gate.log
	LogLevel               string `yaml:"log_level"`                //	log_level: debug
	NetTimeout             string `yaml:"net_timeout"`              //	net_timeout: 1000
	SendSize               string `yaml:"send_size"`                //	send_size:
	ReadSize               string `yaml:"read_size"`                //	read_size:
	RsaKey                 string `yaml:"rsa_key"`                  //	rsa_key: rsa.key
	RsaCertificate         string `yaml:"rsa_certificate"`          //	rsa_certificate: rsa.crt
	CompressConnection     string `yaml:"compress_connection"`      //	compress_connection: false
	EncryptConnection      string `yaml:"encrypt_connection"`       //	encrypt_connection: false
	HeartbeatCheckInterval string `yaml:"heartbeat_check_interval"` //	heartbeat_check_interval: 0

	Net  string `yaml:"net"` //     net: "tcp" # kcp, websocket
	IP   string `yaml:"ip"`
	Port string `yaml:"port"` //     port: 49551

	Mongo mongo.Config `yaml:"mongo"`
}

var CFG Config

func InitConfig() {
	conf.Init("proxy.yaml", &CFG)

	xlog.Infoln("CFG:", CFG)
}
