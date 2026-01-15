package network

import (
	"github.com/0x822a5b87/tiny-docker/src/entity"
)

type IdStore interface {
}

type NetworkStore interface {
	GetAll() ([]*entity.Network, error)
	Get(networkId entity.NetworkId) (*entity.Network, error)
	Update(networkId entity.NetworkId, network *entity.Network) error
	Delete(networkId entity.NetworkId) error
	GetByName(name string) (*entity.Network, error)
}

type EndpointStore interface {
	Get(endpointId entity.EndpointId) (*entity.Endpoint, error)
	Update(endpointId entity.EndpointId, endpoint *entity.Endpoint) error
	Delete(endpointId entity.EndpointId) error
}

type IPAMStore interface {
	Get(networkId entity.NetworkId) (IPAM, error)
	Update(networkId entity.NetworkId, ipam IPAM) error
	Delete(networkId entity.NetworkId) error
}
