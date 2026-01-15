package network

import (
	"net"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

type NetworkDriver interface {
	Create(name string, subnet *net.IPNet) (*entity.Network, error)
	Delete(network *entity.Network) error
	Connect(network *entity.Network, endpoint *entity.Endpoint) error
	Disconnect(network *entity.Network, endpoint *entity.Endpoint) error
}

type BridgeDriver struct{}

func (driver *BridgeDriver) Create(name string, subnet *net.IPNet) (*entity.Network, error) {
	if err := util.CreateBridge(name); err != nil {
		logrus.Errorf("error create bridge : %s, err : %s", name, err)
		return nil, err
	}

	id, _ := conf.GenUUID()
	network := &entity.Network{
		Id:      entity.NetworkId(id),
		Name:    name,
		Type:    entity.NetworkBridge,
		IPNet:   subnet,
		Gateway: subnet.IP,
	}

	if err := util.SetBridgeIP(name, subnet.String()); err != nil {
		logrus.Errorf("error set bridge ip : %s", err)
		return nil, err
	}

	if err := util.SetupLinkByName(name); err != nil {
		logrus.Errorf("error setup bridge ip : %s", err)
		return nil, err
	}

	return network, nil
}

func (driver *BridgeDriver) Delete(network *entity.Network) error {
	//TODO implement me
	panic("implement me")
}

func (driver *BridgeDriver) Connect(network *entity.Network, endpoint *entity.Endpoint) error {
	//TODO implement me
	panic("implement me")
}

func (driver *BridgeDriver) Disconnect(network *entity.Network, endpoint *entity.Endpoint) error {
	//TODO implement me
	panic("implement me")
}

type IPAM interface {
	AllocateIP() (*net.IP, error)
	ReleaseIP(*net.IP)
}

type IPAMImpl struct {
	AvailableIps []*net.IP `json:"available_ips"`
}

func (ipam *IPAMImpl) AllocateIP() (*net.IP, error) {
	if len(ipam.AvailableIps) == 0 {
		return nil, constant.ErrResourcePoolIsEmpty
	}
	ip := ipam.AvailableIps[len(ipam.AvailableIps)-1]
	ipam.AvailableIps = ipam.AvailableIps[:len(ipam.AvailableIps)-1]
	return ip, nil
}

func (ipam *IPAMImpl) ReleaseIP(ip *net.IP) {
	ipam.AvailableIps = append(ipam.AvailableIps, ip)
}
