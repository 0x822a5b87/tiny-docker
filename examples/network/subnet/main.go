package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const CIDR = "172.17.0.0/16"

type SubnetExample struct {
	Gateway         net.IP     `json:"gateway"`
	IPNet           *net.IPNet `json:"ip_net"`
	AvailableIpPool []net.IP   `json:"available_ip_pool"`
}

func ParseSubnetFromCIDR(cidr string) (*SubnetExample, error) {
	subnet := &SubnetExample{}
	var err error
	subnet.Gateway, subnet.IPNet, err = net.ParseCIDR(cidr)
	if err != nil {
		log.Printf("error parsing CIDR: %v", err)
		return nil, err
	}
	subnet.AvailableIpPool, err = GetAllAvailableIp(subnet.IPNet)
	if err != nil {
		log.Printf("error getting available ips: %v", err)
		return nil, err
	}
	return subnet, nil
}

func GetBroadcastIPInt(ipNet *net.IPNet) uint32 {
	networkInt := binary.BigEndian.Uint32(ipNet.IP.To4())
	maskInt := binary.BigEndian.Uint32(ipNet.Mask)
	broadcastInt := networkInt | (^maskInt)
	return broadcastInt
}

func GetNetworkAddress(ipNet *net.IPNet) uint32 {
	networkInt := binary.BigEndian.Uint32(ipNet.IP.To4())
	maskInt := binary.BigEndian.Uint32(ipNet.Mask)
	broadcastInt := networkInt & maskInt
	return broadcastInt
}

func IsValidHostIP(ip net.IP, ipNet *net.IPNet) bool {
	if !ipNet.Contains(ip) {
		return false
	}
	ipInt := binary.BigEndian.Uint32(ip.To4())
	networkInt := GetNetworkAddress(ipNet)
	broadcastInt := GetBroadcastIPInt(ipNet)
	return ipInt != networkInt && ipInt != broadcastInt
}

func IsValidIPv4SubnetCidr(ipNet *net.IPNet) bool {
	if ipNet.IP.To4() == nil {
		return false
	}
	ipInt := binary.BigEndian.Uint32(ipNet.IP.To4())
	networkInt := GetNetworkAddress(ipNet)
	return ipInt == networkInt
}

func GetAllAvailableIp(ipNet *net.IPNet) ([]net.IP, error) {
	if ipNet.IP.To4() == nil {
		return nil, fmt.Errorf("only IPv4 addresses are supported")
	}
	if !IsValidIPv4SubnetCidr(ipNet) {
		return nil, fmt.Errorf("invalid IPv4 CIDR")
	}
	availableIpPool := make([]net.IP, 0)
	networkInt := GetNetworkAddress(ipNet)
	broadcastIPInt := GetBroadcastIPInt(ipNet)
	for availableIp := networkInt + 1; availableIp < broadcastIPInt; availableIp++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, availableIp)
		availableIpPool = append(availableIpPool, ip)
	}
	return availableIpPool, nil
}

func main() {
	subnet, err := ParseSubnetFromCIDR(CIDR)
	if err != nil {
		log.Printf("error parsing CIDR: %v", err)
		return
	}

	firstIp := subnet.AvailableIpPool[0]
	lastIp := subnet.AvailableIpPool[len(subnet.AvailableIpPool)-1]
	log.Printf("gateway = %s, ipNet = %s, first ip = %s, last ip = %s", subnet.Gateway, subnet.IPNet, firstIp, lastIp)
}
