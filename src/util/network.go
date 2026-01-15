package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
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
	parentInt := IpToUint64(parentIP)
	offsetInt := parentInt + nth*subnetSize
	offsetIP := Uint64ToIP(offsetInt)

	targetMask := net.CIDRMask(parentMaskLen, 32)
	targetSubnet := &net.IPNet{
		IP:   offsetIP.To4(),
		Mask: targetMask,
	}

	return targetSubnet, nil
}

func GetNthIp(subnet *net.IPNet, nth int) (*net.IPNet, error) {
	baseIp := subnet.IP.To4()
	if baseIp == nil {
		return nil, constant.ErrNetworkVersion
	}
	ip := IpToUint64(baseIp) + uint64(nth)
	v := &net.IPNet{
		IP:   Uint64ToIP(ip),
		Mask: subnet.Mask,
	}
	return v, nil
}

func GetIpOffset(subnet *net.IPNet, ip *net.IPNet) (int, error) {
	baseIpUint64 := IpToUint64(subnet.IP)
	actualIp := IpToUint64(ip.IP)
	offset := int(actualIp - baseIpUint64)
	return offset, nil
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

	gatewaySubnet, err := GetGateway(subnet)
	if err != nil {
		return "", err
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

func GetGateway(ipNet *net.IPNet) (*net.IPNet, error) {
	networkIP := ipNet.IP.To4()
	if networkIP == nil {
		return nil, constant.ErrNetworkVersion
	}

	gatewayIP := make(net.IP, len(networkIP))
	copy(gatewayIP, networkIP)
	gatewayIP[3] += 1

	gatewayIpNet := &net.IPNet{
		IP:   gatewayIP,
		Mask: ipNet.Mask,
	}

	return gatewayIpNet, nil
}

func IpToUint64(ip net.IP) uint64 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint64(ip[0])<<24 | uint64(ip[1])<<16 | uint64(ip[2])<<8 | uint64(ip[3])
}

func Uint64ToIP(n uint64) net.IP {
	return net.IPv4(
		byte(n>>24),
		byte(n>>16),
		byte(n>>8),
		byte(n),
	)
}

func Connect(vethNs, vethHost, bridgeName string, vethNsIp *net.IPNet, containerPid int) error {
	if err := createVethPair(vethNs, vethHost, containerPid); err != nil {
		logrus.Errorf("[Connect]error create veth pair %v", err)
		return err
	}
	if err := setMaster(vethHost, bridgeName); err != nil {
		logrus.Errorf("[Connect]error set master %v", err)
		return err
	}
	if err := setUpStatus(bridgeName, vethHost, vethNs, containerPid); err != nil {
		logrus.Errorf("error setup status%v\n", err)
		return err
	}
	if err := setVethIPInNS(vethNs, vethNsIp, containerPid); err != nil {
		log.Printf("error setup veth ip %v\n", err)
		return err
	}
	return nil
}

func DeleteVeth(vethName string) error {
	// delete veth host
	link, err := netlink.LinkByName(vethName)
	if err != nil {
		logrus.Errorf("[DeleteVeth]error link by name = %s, err = %v", vethName, err)
		return err
	}
	if err = netlink.LinkDel(link); err != nil {
		logrus.Errorf("[DeleteVeth]error link del : err = %v", err)
		return err
	}
	return nil
}

func createVethPair(vethNs, vethHost string, pid int) error {
	ns, err := netns.GetFromPid(pid)
	if err != nil {
		logrus.Errorf("[createVethPair]error get ns from pid: pid = %d, err = %v", pid, err)
		return err
	}
	la := netlink.LinkAttrs{
		Name:      vethNs,
		Namespace: netlink.NsFd(ns),
	}
	veth := &netlink.Veth{
		LinkAttrs: la,
		PeerName:  vethHost,
	}
	if err = netlink.LinkAdd(veth); err != nil {
		logrus.Errorf("[createVethPair]error link veth pair: %v", err)
		return err
	}
	logrus.Info("Created veth pair: %s and %s\n", vethNs, vethHost)
	return nil
}

func setMaster(vethHostName string, bridgeName string) error {
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		logrus.Errorf("[setMaster]error link bridge by name : %s, err = %v", bridgeName, err)
		return err
	}
	bridge, ok := link.(*netlink.Bridge)
	if !ok {
		logrus.Errorf("[setMaster]error %s is not bridge, err = %v", bridgeName, err)
		return constant.ErrMalformedType
	}
	veth, err := netlink.LinkByName(vethHostName)
	if err != nil {
		logrus.Errorf("[setMaster]error link veth by name : %s, err = %v", vethHostName, err)
		return err
	}
	err = netlink.LinkSetMaster(veth, bridge)
	if err != nil {
		logrus.Errorf("[setMaster]error set master: veth = %s, bridge = %s, err = %v", veth, bridge, err)
		return err
	}

	return nil
}

func setupLinkByName(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		logrus.Errorf("[setupLinkByName]error link by name: name = %s, err = %v", name, err)
		return err
	}
	if err = netlink.LinkSetUp(link); err != nil {
		logrus.Errorf("[setupLinkByName]error set up link: name = %s, err = %v", name, err)
		return err
	}
	return nil
}

func setUpStatus(bridgeName, vethHost, vethNs string, pid int) error {
	if err := setupLinkByName(vethHost); err != nil {
		return err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origins, newNs, err := enterNs(pid)
	if err != nil {
		logrus.Errorf("[setUpStatus]error enter namespace: ns pid = %d, err = %v", pid, err)
		return err
	}

	if err = setupLinkByName(vethNs); err != nil {
		logrus.Errorf("[setUpStatus]error setup link for namespace: ns pid = %d, err = %v", pid, err)
		return err
	}

	if err = exitNs(origins, newNs); err != nil {
		logrus.Errorf("[setUpStatus]error exit namespace: ns pid = %d, err = %v", pid, err)
		return err
	}

	return setupLinkByName(bridgeName)

}

func setVethIPInNS(vethName string, ip *net.IPNet, pid int) error {
	targetNS, err := netns.GetFromPid(pid)
	if err != nil {
		return err
	}
	defer func() { _ = targetNS.Close() }()

	currentNS, err := netns.Get()
	if err != nil {
		return err
	}
	defer func() { _ = currentNS.Close() }()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err = netns.Set(targetNS); err != nil {
		return err
	}
	defer func() { _ = netns.Set(currentNS) }()

	veth, err := netlink.LinkByName(vethName)
	if err != nil {
		return err
	}

	ipAddr := &netlink.Addr{
		IPNet: ip,
		Peer:  nil,
	}

	return netlink.AddrAdd(veth, ipAddr)
}

func enterNs(pid int) (netns.NsHandle, netns.NsHandle, error) {
	origins, err := netns.Get()
	if err != nil {
		logrus.Errorf("[enterNs]error get origin ns: err = %v", err)
		return 0, 0, err
	}

	newNs, err := netns.GetFromPid(pid)
	if err != nil {
		logrus.Errorf("[enterNs]error get ns from pid: pid = %d, err = %v", pid, err)
		return 0, 0, err
	}

	if err = netns.Set(newNs); err != nil {
		logrus.Errorf("[enterNs]error enter ns: err = %v", err)
		return 0, 0, err
	}

	return origins, newNs, nil
}

func exitNs(origins, newNs netns.NsHandle) error {
	defer func() { _ = origins.Close() }()
	defer func() { _ = newNs.Close() }()
	return netns.Set(origins)
}
