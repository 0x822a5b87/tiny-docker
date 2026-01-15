package network

import (
	"errors"
	"net"
	"sync"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func NewNetworks() (*Networks, error) {
	bitmap, err := NewIPNetBitmap(constant.SizeOfSubnet, constant.SizeOfSubnetIp)
	if err != nil {
		return nil, err
	}
	networkStore, err := NewFileNetworkStore()
	if err != nil {
		return nil, err
	}

	ipamStore, err := NewFileIPAMStore()
	if err != nil {
		return nil, err
	}
	return &Networks{
		Mutex:              sync.Mutex{},
		networkStore:       networkStore,
		endpointStore:      NewInMemoryEndpointStore(),
		ipamStore:          ipamStore,
		networkDriver:      &BridgeDriver{},
		bitmap:             bitmap,
		endpointsInNetwork: make(map[entity.NetworkId][]*entity.Endpoint),
	}, nil
}

type Networks struct {
	sync.Mutex
	networkStore       NetworkStore
	endpointStore      EndpointStore
	ipamStore          IPAMStore
	networkDriver      NetworkDriver
	bitmap             *IPNetBitmap
	endpointsInNetwork map[entity.NetworkId][]*entity.Endpoint
}

func (n *Networks) CreateNetwork(networkType entity.NetworkType, networkName string) error {
	_, err := n.networkStore.GetByName(networkName)
	if err == nil {
		logrus.Errorf("network with name %s already exists", networkName)
		return constant.ErrResourceExists
	}

	if IsResourceNotFound(err) {
		n.Lock()
		defer n.Unlock()
		var network *entity.Network
		network, err = n.createNonExitedNetwork(networkType, networkName)
		if err != nil {
			return err
		}
		err = n.networkStore.Update(network.Id, network)
		if err != nil {
			logrus.Errorf("error updating network: %s", err)
			return err
		}
	}

	return err
}

func (n *Networks) GetNetworkByName(networkName string) (*entity.Network, error) {
	return n.networkStore.GetByName(networkName)
}

func (n *Networks) DeleteNetwork(id entity.NetworkId) error {
	nw, err := n.networkStore.Get(id)
	if err != nil {
		logrus.Errorf("[DeleteNetwork]error getting network: %s", err)
		return err
	}
	endpoints := n.getEndpointsOfNetwork(id)
	if len(endpoints) > 0 {
		return constant.ErrDeviceIsBusy
	}

	if err = n.networkStore.Delete(id); err != nil {
		logrus.Errorf("[DeleteNetwork]error deleting network: %s", err)
		return err
	}

	if err = n.ipamStore.Delete(id); err != nil {
		logrus.Errorf("[DeleteNetwork]error deleting network: %s", err)
		return err
	}

	return util.DeleteDevice(nw.Name)
}

func (n *Networks) Connect(networkName string, container entity.Container) error {
	network, err := n.networkStore.GetByName(networkName)
	if err != nil {
		logrus.Errorf("[Connect]error getting network: %s", err)
		return err
	}
	endpoint, err := n.endpointStore.Get(entity.EndpointId(container.Id))
	if err == nil {
		logrus.Errorf("[Connect]endpoint has been created: %s", endpoint.Name)
		return constant.ErrResourceExists
	}
	ipam, err := n.ipamStore.Get(network.Id)
	if err != nil {
		logrus.Errorf("[Connect]error getting ipam: %s", err)
		return err
	}
	ipNet, err := ipam.AllocateIP()
	if err != nil {
		logrus.Errorf("[Connect]error allocating ip: %s", err)
		return err
	}
	vethHost, vethNs := n.getVethInfo(container)
	return util.Connect(vethNs, vethHost, network.Name, ipNet, container.Pid)
}

func (n *Networks) getEndpointsOfNetwork(id entity.NetworkId) []*entity.Endpoint {
	endpoints, ok := n.endpointsInNetwork[id]
	if !ok {
		return make([]*entity.Endpoint, 0)
	}
	return endpoints
}

// Assume that all callers have acquired the lock when calling this function.
func (n *Networks) createNonExitedNetwork(networkType entity.NetworkType, networkName string) (*entity.Network, error) {
	for {
		// NOTE THAT THE GENERATED NETWORK MAY NOT BE AVAILABLE BECAUSE IT IS USED BY OTHER PROCESSES.
		ipBitmap, subnet, err := n.getAvailableIpNet()
		if err != nil {
			logrus.Errorf("Network generation error: %v", err)
			return nil, err
		}
		network, err := n.networkDriver.Create(networkName, subnet)
		if n.isIPOrIPNetBeingUsedErr(err) {
			continue
		}
		ipam, err := NewBitmapIPAM(subnet, ipBitmap)
		if err != nil {
			return nil, err
		}
		err = n.ipamStore.Update(network.Id, ipam)
		return network, err
	}
}

func (n *Networks) getAvailableIpNet() (*Bitmap, *net.IPNet, error) {
	bitmap, subnetPos, err := n.bitmap.AllocateSubnet()
	if err != nil {
		logrus.Errorf("error allocating subnet : %v", err)
		return nil, nil, err
	}

	_, subnet, err := net.ParseCIDR(constant.BaseCidr)
	if err != nil {
		return nil, nil, err
	}
	if !util.IsValidIPv4SubnetCidr(subnet) {
		logrus.Errorf("not valid IPv4 CIDR for : %v", constant.BaseCidr)
		return nil, nil, constant.ErrNetworkVersion
	}
	subnet, err = util.GetNthSubnet(subnet, subnetPos)
	return bitmap, subnet, err
}

func (n *Networks) isIPOrIPNetBeingUsedErr(err error) bool {
	return errors.Is(err, constant.ErrInvalidGateway) || errors.Is(err, constant.ErrInvalidIp)
}

func (n *Networks) getVethInfo(container entity.Container) (vethHost, vethNs string) {
	prefix := "veth-"
	containerId := string(container.Id)
	maxIDLen := 15 - len(prefix)
	var idPart string
	if len(containerId) > maxIDLen {
		idPart = containerId[:maxIDLen]
	} else {
		idPart = containerId
	}
	vethHost = prefix + idPart
	vethNs = "ns-" + idPart
	return
}
