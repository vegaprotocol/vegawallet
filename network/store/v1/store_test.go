package v1_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/network"
	"code.vegaprotocol.io/go-wallet/network/store/v1"
	vgtest "code.vegaprotocol.io/shared/libs/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Saving already existing network succeeds", testFileStoreV1SaveAlreadyExistingNetworkSucceeds)
	t.Run("Saving network succeeds", testFileStoreV1SaveNetworkSucceeds)
	t.Run("Saving legacy network succeeds", testFileStoreV1SaveLegacyNetworkSucceeds)
	t.Run("Verifying non-existing network fails", testFileStoreV1VerifyingNonExistingNetworkFails)
	t.Run("Verifying existing network succeeds", testFileStoreV1VerifyingExistingNetworkSucceeds)
	t.Run("Getting non-existing network fails", testFileStoreV1GetNonExistingNetworkFails)
	t.Run("Getting existing network succeeds", testFileStoreV1GetExistingNetworkSucceeds)
	t.Run("Getting legacy network succeeds", testFileStoreV1GetLegacyNetworkSucceeds)
	t.Run("Getting network path succeeds", testFileStoreV1GetNetworkPathSucceeds)
	t.Run("Getting networks path succeeds", testFileStoreV1GetNetworksPathSucceeds)
	t.Run("Listing networks succeeds", testFileStoreV1ListingNetworksSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	s, err := v1.InitialiseStore(vegaHome.Paths())

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, vegaHome.NetworksHome())
}

func testFileStoreV1SaveAlreadyExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveNetwork(net)

	// then
	require.NoError(t, err)
}

func testFileStoreV1SaveNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1SaveLegacyNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		Name: "test",
		Nodes: network.GRPCConfig{
			Hosts:   []string{"node-1", "node-2"},
			Retries: 5,
		},
		API: network.APIConfig{
			GRPC: network.GRPCConfig{},
		},
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	expectedNet := network.Network{
		Name:  "test",
		Nodes: network.GRPCConfig{},
		API: network.APIConfig{
			GRPC: network.GRPCConfig{
				Hosts:   []string{"node-1", "node-2"},
				Retries: 5,
			},
		},
	}
	assert.Equal(t, expectedNet, *returnedNet)
}

func testFileStoreV1VerifyingNonExistingNetworkFails(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)

	// when
	exists, err := s.NetworkExists("test")

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(net.Name))

	// when
	exists, err := s.NetworkExists("test")

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingNetworkFails(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)

	// when
	keys, err := s.GetNetwork("test")

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1GetLegacyNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	legacyNet := &network.Network{
		Name: "test",
		Nodes: network.GRPCConfig{
			Hosts:   []string{"node-1", "node-2"},
			Retries: 5,
		},
		API: network.APIConfig{
			GRPC: network.GRPCConfig{},
		},
	}

	// when
	err := s.SaveNetwork(legacyNet)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(legacyNet.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	expectedNet := network.Network{
		Name:  "test",
		Nodes: network.GRPCConfig{},
		API: network.APIConfig{
			GRPC: network.GRPCConfig{
				Hosts:   []string{"node-1", "node-2"},
				Retries: 5,
			},
		},
	}
	assert.Equal(t, expectedNet, *returnedNet)
}

func testFileStoreV1GetNetworkPathSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)

	// when
	returnedPath := s.GetNetworkPath("test")

	// then
	assert.Equal(t, vegaHome.NetworkPath("test"), returnedPath)
}

func testFileStoreV1GetNetworksPathSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)

	// when
	returnedPath := s.GetNetworksPath()

	// then
	assert.Equal(t, vegaHome.NetworksHome(), returnedPath)
}

func testFileStoreV1ListingNetworksSucceeds(t *testing.T) {
	vegaHome := newVegaHome()
	defer vegaHome.Remove()

	// given
	s := InitialiseFromPath(vegaHome)
	net := &network.Network{
		// we use "toml" as name on purpose since we want to verify it's not
		// stripped by the ListNetwork() function.
		Name: "toml",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, vegaHome.NetworkPath(net.Name))

	// when
	nets, err := s.ListNetworks()

	// then
	require.NoError(t, err)
	assert.Equal(t, []string{"toml"}, nets)
}

func InitialiseFromPath(h vegaHome) *v1.Store {
	s, err := v1.InitialiseStore(h.Paths())
	if err != nil {
		panic(err)
	}
	return s
}
