package host_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/0x00b/gobbq/tool/host"
)

func TestMain(m *testing.T) {

	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(netInterfaces)

	fmt.Println(host.GetHostName())
	fmt.Println(host.GetMacAddrs())
	fmt.Println(host.GetIPs())

}
