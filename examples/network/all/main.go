package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func CreateBridge(name string) {
	link, _ := netlink.LinkByName(name)
	if link != nil {
		return
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = name
	bridge := &netlink.Bridge{LinkAttrs: attrs}
	if err := netlink.LinkAdd(bridge); err != nil {
		log.Printf("%v\n", err)
	}
}

func DeleteBridge(name string) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		log.Printf("%v\n", err)
	}

	err = netlink.LinkDel(link)
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func CheckBridge(name string) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("Interface: {index = %d, name = %s}", iface.Index, iface.Name)
}

func CreateNs(newNsName string) {
	// Lock the OS Thread so we don't accidentally switch namespaces
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Save the current network namespace
	origins, err := netns.Get()
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer func() { _ = origins.Close() }()

	// Create a new network namespace
	newNs, _ := netns.NewNamed(newNsName)
	defer func() { _ = newNs.Close() }()

	// return to origin ns
	if err = netns.Set(origins); err != nil {
		log.Printf("%v\n", err)
	}
}

func CheckNs(nsName string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ns, err := netns.GetFromName(nsName)
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("name = %s, id = %s, string = %s, open = %v", nsName, ns.UniqueId(), ns.String(), ns.IsOpen())
}

func CreateVethPair(vethNs, vethHost, nsName string) {
	ns, err := netns.GetFromName(nsName)
	if err != nil {
		log.Printf("%v\n", err)
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
		log.Printf("%v\n", err)
	}
	log.Printf("Created veth pair: %s and %s\n", vethNs, vethHost)
}

func DeleteVeth(vethName string) {
	// delete veth host
	link, err := netlink.LinkByName(vethName)
	if err != nil {
		log.Printf("%v\n", err)
	}
	if err = netlink.LinkDel(link); err != nil {
		log.Printf("%v\n", err)
	}
}

func CheckVethPair(vethNs, vethHost, nsName string) {
	link, err := netlink.LinkByName(vethHost)
	if err != nil {
		log.Printf("%v\n", err)
	}
	veth, ok := link.(*netlink.Veth)
	if !ok {
		panic(fmt.Errorf("%s is not veth", vethHost))
	}

	peerIndex, err := netlink.VethPeerIndex(veth)
	if err != nil {
		log.Printf("%v\n", err)
	}

	log.Printf("veth index = %d, pair index = %d", veth.Index, peerIndex)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Save the current network namespace
	origins, err := netns.Get()
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer func() { _ = origins.Close() }()

	// Create a new network namespace
	newNs, _ := netns.GetFromName(nsName)
	defer func() { _ = newNs.Close() }()

	err = netns.Set(newNs)
	if err != nil {
		log.Printf("%v\n", err)
	}

	link, err = netlink.LinkByName(vethNs)
	if err != nil {
		log.Printf("%v\n", err)
	}
	veth, ok = link.(*netlink.Veth)
	if !ok {
		panic(fmt.Errorf("%s is not veth", vethNs))
	}
	peerIndex, err = netlink.VethPeerIndex(veth)
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("veth index = %d, pair index = %d", veth.Index, peerIndex)

	if err = netns.Set(origins); err != nil {
		log.Printf("%v\n", err)
	}
}

func SetMaster(vethName string, bridgeName string) {
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	bridge, ok := link.(*netlink.Bridge)
	if !ok {
		log.Printf("%s is not bridge", bridgeName)
		return
	}
	veth, err := netlink.LinkByName(vethName)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	err = netlink.LinkSetMaster(veth, bridge)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func setupLinkByName(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	if err = netlink.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}

func SetUpStatus(bridgeName, vethHost, vethNs, nsName string) error {
	if err := setupLinkByName(bridgeName); err != nil {
		return err
	}

	if err := setupLinkByName(vethHost); err != nil {
		return err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origins, newNs, err := enterNs(nsName)
	if err != nil {
		return err
	}
	defer func() { _ = exitNs(origins, newNs) }()

	return setupLinkByName(vethNs)
}

func enterNs(nsName string) (netns.NsHandle, netns.NsHandle, error) {
	origins, err := netns.Get()
	if err != nil {
		return 0, 0, err
	}

	newNs, err := netns.GetFromName(nsName)
	if err != nil {
		return 0, 0, err
	}

	if err = netns.Set(newNs); err != nil {
		return 0, 0, err
	}

	return origins, newNs, nil
}

func exitNs(origins, newNs netns.NsHandle) error {
	defer func() { _ = origins.Close() }()
	defer func() { _ = newNs.Close() }()
	return netns.Set(origins)
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

func SetVethIPInNS(nsName, vethName, ipWithCIDR string) error {
	targetNS, err := netns.GetFromName(nsName)
	if err != nil {
		return err
	}
	defer targetNS.Close()

	currentNS, err := netns.Get()
	if err != nil {
		return err
	}
	defer currentNS.Close()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := netns.Set(targetNS); err != nil {
		return err
	}
	defer netns.Set(currentNS)

	veth, err := netlink.LinkByName(vethName)
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

	_ = netlink.AddrDel(veth, ipAddr)
	return netlink.AddrAdd(veth, ipAddr)
}

func main() {
	bridge := "my-bridge"
	ns := "my-ns"
	vethHost := "my-veth-host"
	vethNs := "my-veth-ns"
	bridgeIP := "10.200.0.1/16"
	vethNsIP := "10.200.0.2/16"

	CreateBridge(bridge)
	defer DeleteBridge(bridge)
	CheckBridge(bridge)

	CreateNs(ns)
	CheckNs(ns)

	CreateVethPair(vethNs, vethHost, ns)
	// A veth pair will be removed once either end of it is removed.
	defer DeleteVeth(vethHost)
	CheckVethPair(vethNs, vethHost, ns)

	SetMaster(vethHost, bridge)
	if err := SetUpStatus(bridge, vethHost, vethNs, ns); err != nil {
		log.Printf("error setup status%v\n", err)
	}

	if err := SetBridgeIP(bridge, bridgeIP); err != nil {
		log.Printf("error setup bridge ip %v\n", err)
	}

	if err := SetVethIPInNS(ns, vethNs, vethNsIP); err != nil {
		log.Printf("error setup veth ip %v\n", err)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
