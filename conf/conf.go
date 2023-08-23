package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type NetConf struct {
	Net  string `yaml:"net"` //     net: "tcp" # kcp, websocket
	IP   string `yaml:"ip"`
	Port string `yaml:"port"` //     port: 49551
}

func Init(confFile string, cfg any) {
	err := ParseYamlConf(confFile, cfg)
	if err != nil {
		panic(err)
	}
}

func ParseYamlConf(fpath string, cfg any) error {

	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}
