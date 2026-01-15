package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func DeleteDevice(name string) error {
	link, _ := netlink.LinkByName(name)
	if link == nil {
		logrus.Errorf("Bridge %s not exists", name)
		return constant.ErrResourceNotExists
	}
	if err := netlink.LinkDel(link); err != nil {
		logrus.Errorf("error delete device: %s, error : %s", name, err.Error())
	}
	return nil
}

func CreateBridge(name string) error {
	link, _ := netlink.LinkByName(name)
	if link != nil {
		logrus.Errorf("Bridge %s already exists", name)
		return constant.ErrResourceExists
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = name
	bridge := &netlink.Bridge{LinkAttrs: attrs}
	if err := netlink.LinkAdd(bridge); err != nil {
		logrus.Errorf("error create bridge : %s", name)
		return err
	}
	return nil
}

func SetBridgeIP(bridgeName, ipWithCIDR string) error {
	bridge, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	ipNet, err := netlink.ParseIPNet(ipWithCIDR)
	if err != nil {
		return err
	}
	ipAddr := &netlink.Addr{
		IPNet: ipNet,
		Peer:  nil,
	}

	_ = netlink.AddrDel(bridge, ipAddr)
	return netlink.AddrAdd(bridge, ipAddr)
}

func SetupLinkByName(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	if err = netlink.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}

func GetBroadcastIPUint32(ipNet *net.IPNet) uint32 {
	networkInt := binary.BigEndian.Uint32(ipNet.IP.To4())
	maskInt := binary.BigEndian.Uint32(ipNet.Mask)
	broadcastInt := networkInt | (^maskInt)
	return broadcastInt
}

func GetNetworkAddressUint32(ipNet *net.IPNet) uint32 {
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
	networkInt := GetNetworkAddressUint32(ipNet)
	broadcastInt := GetBroadcastIPUint32(ipNet)
	return ipInt != networkInt && ipInt != broadcastInt
}

func IsValidIPv4SubnetCidr(ipNet *net.IPNet) bool {
	if ipNet.IP.To4() == nil {
		return false
	}
	ipInt := binary.BigEndian.Uint32(ipNet.IP.To4())
	networkInt := GetNetworkAddressUint32(ipNet)
	return ipInt == networkInt
}

func GetNthSubnet(subnet *net.IPNet, nth uint64) (*net.IPNet, error) {
	if subnet == nil {
		return nil, errors.New("subnet is nil")
	}
	parentIP := subnet.IP.To4()
	if parentIP == nil {
		return nil, errors.New("only support IPv4 subnet")
	}
	parentMaskLen, _ := subnet.Mask.Size()
	if parentMaskLen < 0 || parentMaskLen > 32 {
		return nil, errors.New("invalid parent subnet mask length")
	}

	subnetSize := uint64(1) << (32 - uint64(parentMaskLen))
	parentTotalIP := uint64(1) << (32 - uint64(parentMaskLen))
	totalSubnets := parentTotalIP / subnetSize
	if nth >= totalSubnets {
		return nil, fmt.Errorf("nth exceeds total available subnets: nth = {%d}, total = {%d}", nth, totalSubnets)
	}

	parentInt := ipToUint32(parentIP)
	offsetInt := parentInt + nth*subnetSize
	offsetIP := uint32ToIP(offsetInt)

	targetMask := net.CIDRMask(parentMaskLen, 32)
	targetSubnet := &net.IPNet{
		IP:   offsetIP.To4(),
		Mask: targetMask,
	}

	if !subnet.Contains(targetSubnet.IP) {
		return nil, errors.New("offset subnet is out of parent subnet range")
	}

	return targetSubnet, nil
}

func ipToUint32(ip net.IP) uint64 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint64(ip[0])<<24 | uint64(ip[1])<<16 | uint64(ip[2])<<8 | uint64(ip[3])
}

func uint32ToIP(n uint64) net.IP {
	return net.IPv4(
		byte(n>>24),
		byte(n>>16),
		byte(n>>8),
		byte(n),
	)
}
