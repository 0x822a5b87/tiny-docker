package entity

import "net"

const NetworkBridge NetworkType = "bridge"

type NetworkType string

type NetworkId string

type EndpointId string

type Endpoint struct {
	Id      EndpointId `json:"id"`
	Name    string     `json:"name"`
	MAC     string     `json:"mac"`
	IP      *net.IP    `json:"ip"`
	Network *Network   `json:"network"`
}

type Network struct {
	Id      NetworkId   `json:"id"`
	Name    string      `json:"name"`
	Type    NetworkType `json:"type"`
	Gateway net.IP      `json:"gateway"`
	IPNet   *net.IPNet  `json:"ip_net"`
}
