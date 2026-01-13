package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func main() {
	la := netlink.NewLinkAttrs()
	la.Name = "foo"
	myBridge := &netlink.Bridge{LinkAttrs: la}
	err := netlink.LinkAdd(myBridge)
	if err != nil {
		fmt.Printf("could not add %s: %v\n", la.Name, err)
	}
	eth1, err := netlink.LinkByName("eth1")
	if err != nil {
		panic(err)
	}
	if err = netlink.LinkSetMaster(eth1, myBridge); err != nil {
		panic(err)
	}
}
