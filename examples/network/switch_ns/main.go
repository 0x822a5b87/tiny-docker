package main

import (
	"fmt"
	"net"
	"runtime"

	"github.com/vishvananda/netns"
)

func main() {
	// Lock the OS Thread so we don't accidentally switch namespaces
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Save the current network namespace
	origins, _ := netns.Get()
	defer func() { _ = origins.Close() }()

	// Create a new network namespace
	newNs, _ := netns.New()
	defer func() { _ = newNs.Close() }()

	// Do something with the network namespace
	ifaces, _ := net.Interfaces()
	for _, v := range ifaces {
		fmt.Printf("Interfaces of new namespace: name = %v\n", v.Name)
	}

	// Switch back to the original namespace
	_ = netns.Set(origins)

	ifaces, _ = net.Interfaces()
	for _, v := range ifaces {
		fmt.Printf("Interfaces of host: name = %v\n", v.Name)
	}
}
