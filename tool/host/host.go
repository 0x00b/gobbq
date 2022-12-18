package host

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
)

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetIPs() (ips []string) {

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func GetHostName() string {
	hostname, err1 := os.Hostname()
	if err1 != nil {
		var sum [3]byte
		id := sum[:]
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return string(id)
	}
	return hostname
}

// func main() {
// 	fmt.Printf("mac addrs: %q\n", getMacAddrs())
// 	fmt.Printf("ips: %q\n", getIPs())
// 	hostname, _ := os.Hostname()
// 	fmt.Printf("ips: %q\n", hostname)
// }
