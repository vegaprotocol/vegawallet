package v1_test

import (
	"os"
	"path/filepath"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	vgtest "code.vegaprotocol.io/shared/libs/test"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/network"
	v1 "code.vegaprotocol.io/vegawallet/network/store/v1"
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
	t.Run("Getting network path succeeds", testFileStoreV1GetNetworkPathSucceeds)
	t.Run("Getting networks path succeeds", testFileStoreV1GetNetworksPathSucceeds)
	t.Run("Listing networks succeeds", testFileStoreV1ListingNetworksSucceeds)
	t.Run("Deleting network succeeds", testFileStoreV1DeleteNetworkSucceeds)
}

func testFileStoreV1DeleteNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// Create a network for us to delete
	s, err := v1.InitialiseStore(vegaHome)
	require.NoError(t, err)
	assert.NotNil(t, s)

	net := &network.Network{
		Name: "test",
	}

	err = s.SaveNetwork(net)
	require.NoError(t, err)

	// Check it's really there
	returnedNet, err := s.GetNetwork("test")
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)

	// Now delete it
	err = s.DeleteNetwork("test")
	require.NoError(t, err)

	// Check it's no longer there
	returnedNet, err = s.GetNetwork("test")
	require.Error(t, err)
	assert.Nil(t, returnedNet)
}

func testNewStoreSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	s, err := v1.InitialiseStore(vegaHome)

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, networksHome(t, vegaHome))
}

func testFileStoreV1SaveAlreadyExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
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
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, networkPath(t, vegaHome, net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1VerifyingNonExistingNetworkFails(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	exists, err := s.NetworkExists("test")

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, networkPath(t, vegaHome, net.Name))

	// when
	exists, err := s.NetworkExists("test")

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingNetworkFails(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	keys, err := s.GetNetwork("test")

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingNetworkSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	net := &network.Network{
		Name: "test",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, networkPath(t, vegaHome, net.Name))

	// when
	returnedNet, err := s.GetNetwork("test")

	// then
	require.NoError(t, err)
	assert.Equal(t, net, returnedNet)
}

func testFileStoreV1GetNetworkPathSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	returnedPath := s.GetNetworkPath("test")

	// then
	assert.Equal(t, networkPath(t, vegaHome, "test"), returnedPath)
}

func testFileStoreV1GetNetworksPathSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	returnedPath := s.GetNetworksPath()

	// then
	assert.Equal(t, networksHome(t, vegaHome), returnedPath)
}

func testFileStoreV1ListingNetworksSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	net := &network.Network{
		// we use "toml" as name on purpose since we want to verify it's not
		// stripped by the ListNetwork() function.
		Name: "toml",
	}

	// when
	err := s.SaveNetwork(net)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, networkPath(t, vegaHome, net.Name))

	// when
	nets, err := s.ListNetworks()

	// then
	require.NoError(t, err)
	assert.Equal(t, []string{"toml"}, nets)
}

func initialiseFromPath(t *testing.T, vegaHome paths.Paths) *v1.Store {
	t.Helper()
	s, err := v1.InitialiseStore(vegaHome)
	if err != nil {
		t.Fatalf("couldn't initialise store: %v", err)
	}
	return s
}

func newVegaHome(t *testing.T) *paths.CustomPaths {
	t.Helper()
	rootPath := filepath.Join("/tmp", "vegawallet", vgrand.RandomStr(10))
	t.Cleanup(func() {
		if err := os.RemoveAll(rootPath); err != nil {
			t.Fatalf("couldn't remove vega home: %v", err)
		}
	})

	return &paths.CustomPaths{CustomHome: rootPath}
}

func networksHome(t *testing.T, vegaHome *paths.CustomPaths) string {
	t.Helper()
	return vegaHome.ConfigPathFor(paths.WalletServiceNetworksConfigHome)
}

func networkPath(t *testing.T, vegaHome *paths.CustomPaths, name string) string {
	t.Helper()
	return filepath.Join(networksHome(t, vegaHome), name+".toml")
}
