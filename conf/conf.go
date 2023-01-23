package conf

type Conf struct {
	Game  GameConf  `yaml:"game"`
	Gate  GateConf  `yaml:"gate"`
	Proxy ProxyConf `yaml:"proxy"`
}

type GameConf struct {
	LogFile  string `yaml:"log_file"`  // log_file: gate.log
	LogLevel string `yaml:"log_level"` // log_level: debug
}

type ProxyConf struct {
	// common:
	//	log_file: gate.log
	//	log_level: debug
	//	net_timeout: 1000
	//	send_size:
	//	read_size:
	//	rsa_key: rsa.key
	//	rsa_certificate: rsa.crt
	//	compress_connection: false
	//	encrypt_connection: false
	//	heartbeat_check_interval: 0
	//
	// inst:
	//   - id: 1
	//     net: "tcp" # kcp, websocket
	//     port: 49551
}

type GateConf struct {
	//		common:
	//	    log_file: gate.log
	//	    log_level: debug
	//	    net_timeout: 1000
	//	    send_size:
	//	    read_size:
	//	    rsa_key: rsa.key
	//	    rsa_certificate: rsa.crt
	//	    compress_connection: false
	//	    encrypt_connection: false
	//	    heartbeat_check_interval: 0
	//	  inst:
	//	    - id: 1
	//	      net: "tcp" # kcp, websocket
	//	      port: 59551
	//	    - id: 2
	//	      net: "tcp" # kcp, websocket
	//	      port: 59552
}
