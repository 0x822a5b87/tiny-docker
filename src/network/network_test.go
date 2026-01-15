package network

import (
	"testing"

	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/stretchr/testify/assert"
)

func TestNetworks(t *testing.T) {
	util.InitTestConfig()
	networks, err := NewNetworks()
	assert.NoError(t, err)
	assert.NotNil(t, networks)
	bridgeName := "test-bridge"
	err = networks.CreateNetwork(entity.NetworkBridge, bridgeName)
	assert.NoError(t, err)

	container := entity.Container{
		Id:   "c6d79b2fb313442d85f0ff6766a2ea70",
		Name: "test-container",
		Pid:  1198966,
	}
	err = networks.Connect(bridgeName, container)
	assert.NoError(t, err)

	network, err := networks.networkStore.GetByName(bridgeName)
	assert.NoError(t, err)
	assert.NotNil(t, network)
	assert.Equal(t, bridgeName, network.Name)
	err = networks.DeleteNetwork(network.Id)
	assert.NoError(t, err)
}
