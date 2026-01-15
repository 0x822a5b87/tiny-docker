package network

import (
	"net"
	"sync"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
)

func NewInMemoryNetworkStore() NetworkStore {
	return &InMemoryNetworkStore{
		mutex:    sync.RWMutex{},
		names:    make(map[string]entity.NetworkId),
		networks: make(map[entity.NetworkId]*entity.Network),
		pool:     NewRecyclePool[net.IPNet](1024),
	}
}

type InMemoryNetworkStore struct {
	mutex    sync.RWMutex
	names    map[string]entity.NetworkId
	networks map[entity.NetworkId]*entity.Network
	pool     *RecyclePool[net.IPNet]
}

func (store *InMemoryNetworkStore) GetAll() ([]*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	networks := make([]*entity.Network, 0, len(store.names))
	for _, name := range store.names {
		networks = append(networks, store.networks[name])
	}
	return networks, nil
}

func (store *InMemoryNetworkStore) GetByName(name string) (*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	id, ok := store.names[name]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return store.Get(id)
}

func (store *InMemoryNetworkStore) Get(networkId entity.NetworkId) (*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	network, ok := store.networks[networkId]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return network, nil
}

func (store *InMemoryNetworkStore) Update(networkId entity.NetworkId, network *entity.Network) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.networks[networkId] = network
	store.names[network.Name] = network.Id
	return nil
}

func (store *InMemoryNetworkStore) Delete(networkId entity.NetworkId) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.networks, networkId)
	return nil
}

func NewInMemoryEndpointStore() EndpointStore {
	return &InMemoryEndpointStore{
		mutex:     sync.RWMutex{},
		endpoints: make(map[entity.EndpointId]*entity.Endpoint),
	}
}

type InMemoryEndpointStore struct {
	mutex     sync.RWMutex
	endpoints map[entity.EndpointId]*entity.Endpoint
}

func (store *InMemoryEndpointStore) Get(endpointId entity.EndpointId) (*entity.Endpoint, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	endpoint, ok := store.endpoints[endpointId]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return endpoint, nil

}

func (store *InMemoryEndpointStore) Update(endpointId entity.EndpointId, endpoint *entity.Endpoint) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.endpoints[endpointId] = endpoint
	return nil
}

func (store *InMemoryEndpointStore) Delete(endpointId entity.EndpointId) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.endpoints, endpointId)
	return nil
}

func NewInMemoryIPAMStore() IPAMStore {
	return &InMemoryIPAMStore{
		mutex: sync.RWMutex{},
		ipams: make(map[entity.NetworkId]IPAM, 0),
	}
}

type InMemoryIPAMStore struct {
	mutex sync.RWMutex
	ipams map[entity.NetworkId]IPAM
}

func (store *InMemoryIPAMStore) Get(networkId entity.NetworkId) (IPAM, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	ipam, ok := store.ipams[networkId]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return ipam, nil
}

func (store *InMemoryIPAMStore) Update(networkId entity.NetworkId, ipam IPAM) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.ipams[networkId] = ipam
	return nil
}

func (store *InMemoryIPAMStore) Delete(networkId entity.NetworkId) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.ipams, networkId)
	return nil
}
