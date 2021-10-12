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
	t.Run("Verifying non-existing network fails", testFileStoreV1VerifyingNonExistingNetworkFails)
	t.Run("Verifying existing network succeeds", testFileStoreV1VerifyingExistingNetworkSucceeds)
	t.Run("Getting non-existing network fails", testFileStoreV1GetNonExistingNetworkFails)
	t.Run("Getting existing network succeeds", testFileStoreV1GetExistingNetworkSucceeds)
	t.Run("Listing networks succeeds", testFileStoreV1ListingNetworksSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	s, err := v1.InitialiseStore(configDir.Paths())

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, configDir.NetworksHome())
}

func testFileStoreV1SaveAlreadyExistingNetworkSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
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
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, configDir.NetworkPath(net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1VerifyingNonExistingNetworkFails(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	exists, err := s.NetworkExists("test")

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingNetworkSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, configDir.NetworkPath(net.Name))

	// when
	exists, err := s.NetworkExists("test")

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingNetworkFails(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	keys, err := s.GetNetwork("test")

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingNetworkSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, configDir.NetworkPath(net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1ListingNetworksSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	net := &network.Network{
		// we use "toml" as name on purpose since we want to verify it's not
		// stripped by the ListNetwork() function.
		Name: "toml",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, configDir.NetworkPath(net.Name))

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
