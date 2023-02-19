package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var C Conf

type Conf struct {
	Game  GameConf  `yaml:"game"`
	Gate  GateConf  `yaml:"gate"`
	Proxy ProxyConf `yaml:"proxy"`
}

type GameConf struct {
	LogFile  string `yaml:"log_file"`  // log_file: gate.log
	LogLevel string `yaml:"log_level"` // log_level: debug
}

type CommConf struct {
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
}

type Inst struct {
	CommConf
	ID   string `yaml:"id"`  //   - id: 1
	Net  string `yaml:"net"` //     net: "tcp" # kcp, websocket
	IP   string `yaml:"ip"`
	Port string `yaml:"port"` //     port: 49551
}

type ProxyConf struct {
	InstNum uint32   `yaml:"inst_num"`
	Comm    CommConf `yaml:"common"`
	Inst    []Inst   `yaml:"inst"`
}

type GateConf struct {
	InstNum uint32   `yaml:"inst_num"`
	Comm    CommConf `yaml:"common"`
	Inst    []Inst   `yaml:"inst"`
}

func Init(confFile string) {
	err := ParseYamlConf(confFile, &C)
	if err != nil {
		panic(err)
	}
}

func ParseYamlConf(fpath string, cfg any) error {

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}
