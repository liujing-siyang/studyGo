package service

import (
	"fmt"
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	ip,_ := GetOutboundIP()
	fmt.Println(ip.To4())
}

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
