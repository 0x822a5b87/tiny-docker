package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/sirupsen/logrus"
)

func NetworkCreate(name string) (*entity.Network, error) {
	conf.LoadBasicCommand()
	if err := networks.CreateNetwork(entity.NetworkBridge, name); err != nil {
		return nil, err
	}
	return networks.GetNetworkByName(name)
}

func NetworkRm(name string) (*entity.Network, error) {
	conf.LoadBasicCommand()
	network, err := networks.GetNetworkByName(name)
	if err != nil {
		logrus.Error("network rm error: ", err)
		return nil, err
	}
	err = networks.DeleteNetwork(network.Id)
	return network, err
}

func NetworkInspect(name string) (*entity.Network, error) {
	conf.LoadBasicCommand()
	return networks.GetNetworkByName(name)
}
