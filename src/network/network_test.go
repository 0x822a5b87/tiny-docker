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

	network, err := networks.networkStore.GetByName(bridgeName)
	assert.NoError(t, err)
	assert.NotNil(t, network)
	assert.Equal(t, bridgeName, network.Name)
	err = networks.DeleteNetwork(network.Id)
	assert.NoError(t, err)
}
