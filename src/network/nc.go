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
	if _, err := util.CheckSubnetGatewayAvailable(subnet); err != nil {
		logrus.Errorf("subnet gateway not available: %v", err)
		return nil, err
	}

	if err := util.CreateBridge(name); err != nil {
		logrus.Errorf("error create bridge : %s, err : %s", name, err)
		return nil, err
	}

	gateway, err := util.GetGateway(subnet)
	if err != nil {
		return nil, err
	}

	id, _ := conf.GenUUID()
	network := &entity.Network{
		Id:      entity.NetworkId(id),
		Name:    name,
		Type:    entity.NetworkBridge,
		IPNet:   subnet,
		Gateway: gateway,
	}

	gateway, err = util.GetGateway(subnet)
	if err != nil {
		logrus.Errorf("error get gateway : %s, err : %s", gateway, err)
		return nil, err
	}

	if err = util.SetBridgeIP(name, gateway.String()); err != nil {
		logrus.Errorf("error set bridge ip : %s", err)
		return nil, err
	}

	if err = util.SetupLinkByName(name); err != nil {
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
	AllocateIP() (*net.IPNet, error)
	ReleaseIP(*net.IPNet) error
}

func NewBitmapIPAM(subnet *net.IPNet, bitmap *Bitmap) (*BitmapIPAM, error) {
	return &BitmapIPAM{
		Subnet: subnet,
		Bitmap: bitmap,
	}, nil
}

type BitmapIPAM struct {
	Subnet *net.IPNet `json:"subnet"`
	Bitmap *Bitmap    `json:"bitmap"`
}

func (ipam *BitmapIPAM) AllocateIP() (*net.IPNet, error) {
	pos := ipam.Bitmap.FindFirstUnset()
	if pos < 0 {
		return nil, constant.ErrResourcePoolIsEmpty
	}
	return util.GetNthIp(ipam.Subnet, pos)
}

func (ipam *BitmapIPAM) ReleaseIP(ip *net.IPNet) error {
	pos, err := util.GetIpOffset(ipam.Subnet, ip)
	if err != nil {
		return err
	}
	return ipam.Bitmap.Clear(uint64(pos))
}
