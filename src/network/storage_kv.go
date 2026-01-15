package network

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func NewFileNetworkStore() (NetworkStore, error) {
	store := &FileNetworkStore{
		mutex:    sync.RWMutex{},
		names:    make(map[string]entity.NetworkId),
		networks: make(map[entity.NetworkId]*entity.Network),
	}
	networkBasePath := conf.RuntimeNetworkPath.Get()
	if err := util.EnsureFilePathExist(networkBasePath); err != nil {
		return nil, err
	}
	filesInDir, err := util.ReadAllFilesInDir(networkBasePath)
	if err != nil {
		logrus.Errorf("ReadAllFilesInDir error: %s", err)
		return nil, err
	}
	for _, data := range filesInDir {
		network := &entity.Network{}
		if err = json.Unmarshal(data, network); err != nil {
			logrus.Errorf("error unmarshal network from file: %s", err)
			continue
		}
		store.names[network.Name] = network.Id
		store.networks[network.Id] = network
	}

	return store, nil
}

type FileNetworkStore struct {
	mutex    sync.RWMutex
	names    map[string]entity.NetworkId
	networks map[entity.NetworkId]*entity.Network
}

func (store *FileNetworkStore) GetAll() ([]*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	networks := make([]*entity.Network, 0, len(store.names))
	for _, name := range store.names {
		networks = append(networks, store.networks[name])
	}
	return networks, nil
}

func (store *FileNetworkStore) GetByName(name string) (*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	id, ok := store.names[name]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return store.Get(id)
}

func (store *FileNetworkStore) Get(networkId entity.NetworkId) (*entity.Network, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	return store.doGet(networkId)
}

func (store *FileNetworkStore) Update(networkId entity.NetworkId, network *entity.Network) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.networks[networkId] = network
	store.names[network.Name] = network.Id
	data, err := json.Marshal(network)
	if err != nil {
		return err
	}
	if err = os.WriteFile(store.getNetworkFile(networkId), data, 0644); err != nil {
		return err
	}
	return nil
}

func (store *FileNetworkStore) Delete(networkId entity.NetworkId) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	network, err := store.doGet(networkId)
	if err != nil {
		return err
	}
	if err = os.Remove(store.getNetworkFile(networkId)); err != nil {
		return err
	}
	delete(store.networks, networkId)
	delete(store.names, network.Name)
	return nil
}

func (store *FileNetworkStore) getNetworkFile(id entity.NetworkId) string {
	networkBasePath := conf.RuntimeNetworkPath.Get()
	return filepath.Join(networkBasePath, string(id))
}

func (store *FileNetworkStore) doGet(id entity.NetworkId) (*entity.Network, error) {
	network, ok := store.networks[id]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return network, nil
}

func NewFileIPAMStore() (IPAMStore, error) {
	store := &FileIPAMStore{
		mutex:   sync.RWMutex{},
		ipamMap: make(map[entity.NetworkId]IPAM),
	}

	ipamBasePath := conf.RuntimeIpamPath.Get()
	if err := util.EnsureFilePathExist(ipamBasePath); err != nil {
		return nil, err
	}

	filesInDir, err := util.ReadAllFilesInDir(ipamBasePath)
	if err != nil {
		logrus.Errorf("[NewFileIPAMStore]ReadAllFilesInDir error: %s", err)
		return nil, err
	}
	for networkId, data := range filesInDir {
		ipam := &BitmapIPAM{}
		if err = json.Unmarshal(data, ipam); err != nil {
			logrus.Errorf("[NewFileIPAMStore]error unmarshal ipam from file: %s", err)
			continue
		}
		store.ipamMap[entity.NetworkId(networkId)] = ipam
	}

	return store, nil
}

type FileIPAMStore struct {
	mutex   sync.RWMutex
	ipamMap map[entity.NetworkId]IPAM
}

func (store *FileIPAMStore) Get(networkId entity.NetworkId) (IPAM, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	ipam, ok := store.ipamMap[networkId]
	if !ok {
		return nil, constant.ErrResourceNotFound
	}
	return ipam, nil
}

func (store *FileIPAMStore) Update(networkId entity.NetworkId, ipam IPAM) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	data, err := json.Marshal(ipam)
	if err != nil {
		return err
	}
	if err = os.WriteFile(store.getIpamFile(networkId), data, 0644); err != nil {
		return err
	}
	store.ipamMap[networkId] = ipam
	return nil
}

func (store *FileIPAMStore) Delete(networkId entity.NetworkId) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := os.Remove(store.getIpamFile(networkId)); err != nil {
		return err
	}
	delete(store.ipamMap, networkId)
	return nil
}

func (store *FileIPAMStore) getIpamFile(networkId entity.NetworkId) string {
	networkBasePath := conf.RuntimeIpamPath.Get()
	return filepath.Join(networkBasePath, string(networkId))
}
