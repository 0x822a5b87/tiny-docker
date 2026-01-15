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
	return &Networks{
		Mutex:              sync.Mutex{},
		networkStore:       networkStore,
		endpointStore:      NewInMemoryEndpointStore(),
		ipamStore:          NewInMemoryIPAMStore(),
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
		logrus.Errorf("error getting network: %s", err)
		return err
	}
	endpoints := n.getEndpointsOfNetwork(id)
	if len(endpoints) > 0 {
		return constant.ErrDeviceIsBusy
	}

	if err = n.networkStore.Delete(id); err != nil {
		logrus.Errorf("error deleting network: %s", err)
		return err
	}

	return util.DeleteDevice(nw.Name)
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
		subnet, err := n.getAvailableIpNet()
		if err != nil {
			logrus.Errorf("Network generation error: %v", err)
			return nil, err
		}
		network, err := n.networkDriver.Create(networkName, subnet)
		if n.isIPOrIPNetBeingUsedErr(err) {
			continue
		}
		return network, err
	}
}

func (n *Networks) getAvailableIpNet() (*net.IPNet, error) {
	_, subnetPos, err := n.bitmap.AllocateSubnet()
	if err != nil {
		logrus.Errorf("error allocating subnet : %v", err)
		return nil, err
	}

	_, subnet, err := net.ParseCIDR(constant.BaseCidr)
	if err != nil {
		return nil, err
	}
	if !util.IsValidIPv4SubnetCidr(subnet) {
		logrus.Errorf("not valid IPv4 CIDR for : %v", constant.BaseCidr)
		return nil, constant.ErrNetworkVersion
	}
	return util.GetNthSubnet(subnet, subnetPos)
}

func (n *Networks) isIPOrIPNetBeingUsedErr(err error) bool {
	return errors.Is(err, constant.ErrInvalidGateway) || errors.Is(err, constant.ErrInvalidIp)
}
