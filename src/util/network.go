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

	gatewayIp, err := GetGatewayIp(ipNet)
	if err != nil {
		return err
	}
	ipAddr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   gatewayIp,
			Mask: ipNet.Mask,
		},
		Peer: nil,
	}

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

	subnetSize := uint64(1) << uint64(parentMaskLen)
	parentInt := ipToUint32(parentIP)
	offsetInt := parentInt + nth*subnetSize
	offsetIP := uint32ToIP(offsetInt)

	targetMask := net.CIDRMask(parentMaskLen, 32)
	targetSubnet := &net.IPNet{
		IP:   offsetIP.To4(),
		Mask: targetMask,
	}

	return targetSubnet, nil
}

func CheckIPAllocated(targetIP string) (bool, error) {
	ipStr, _, err := net.ParseCIDR(targetIP)
	if err != nil {
		ipStr = net.ParseIP(targetIP)
		if ipStr == nil {
			return false, fmt.Errorf("invalid IP/CIDR format: %s, err: %w", targetIP, err)
		}
	}
	targetIPv4 := ipStr.To4()
	if targetIPv4 == nil {
		return false, errors.New("only support IPv4 address")
	}

	links, err := netlink.LinkList()
	if err != nil {
		return false, fmt.Errorf("netlink list links failed: %w", err)
	}

	for _, link := range links {
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if addr.IP.To4() == nil {
				continue
			}
			if addr.IP.Equal(targetIPv4) {
				logrus.Errorf("IP %s is allocated to device: %s (state: %s)", targetIPv4.String(), link.Attrs().Name, getLinkState(link))
				return true, constant.ErrInvalidIp
			}
		}
	}

	return false, nil
}

func CheckSubnetGatewayAvailable(subnet *net.IPNet) (string, error) {
	networkIP := subnet.IP.To4()
	if networkIP == nil {
		return "", constant.ErrNetworkVersion
	}

	gatewayIp, err := GetGatewayIp(subnet)
	if err != nil {
		return "", err
	}

	gatewaySubnet := &net.IPNet{
		IP:   gatewayIp,
		Mask: subnet.Mask,
	}
	gatewayIpStr := gatewaySubnet.String()

	used, err := CheckIPAllocated(gatewayIpStr)
	if err != nil {
		return "", fmt.Errorf("check gateway IP failed: %w", err)
	}
	if used {
		logrus.Errorf("gateway IP %s is already used by other device", gatewayIpStr)
		return "", constant.ErrInvalidGateway
	}

	return gatewayIpStr, nil
}

func getLinkState(link netlink.Link) string {
	if link.Attrs().Flags&net.FlagUp != 0 {
		return "UP"
	}
	return "DOWN"
}

func GetGatewayIp(ipNet *net.IPNet) ([]byte, error) {
	networkIP := ipNet.IP.To4()
	if networkIP == nil {
		return nil, constant.ErrNetworkVersion
	}

	gatewayIP := make(net.IP, len(networkIP))
	copy(gatewayIP, networkIP)
	gatewayIP[3] += 1

	return gatewayIP, nil
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
