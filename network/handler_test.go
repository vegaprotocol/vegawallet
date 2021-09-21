package network_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/network"
	"code.vegaprotocol.io/go-wallet/network/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testConfig struct {
	ctrl  *gomock.Controller
	log   *zap.Logger
	store *mocks.MockStore
}

func getTestConfig(t *testing.T) *testConfig {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	return &testConfig{
		ctrl:  ctrl,
		log:   zap.NewNop(),
		store: store,
	}
}

func TestGenerateConfig(t *testing.T) {
	t.Run("Initialising config succeeds", testInitialisingConfigSucceeds)
	t.Run("Initialising config with error fails", testInitialisingConfigWithErrorFails)
	t.Run("Initialising config with existing config fails", testInitialisingConfigWithExistingConfigFails)
	t.Run("Overwriting config succeeds", testOverwritingConfigSucceeds)
}

func testInitialisingConfigSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		NetworkExists("fairground").
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveNetwork(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := network.InitialiseNetworks(ts.store, false)

	// then
	require.NoError(t, err)
}

func testInitialisingConfigWithErrorFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		NetworkExists("fairground").
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveNetwork(gomock.Any()).
		Times(1).
		Return(errors.New("some error"))

	// when
	err := network.InitialiseNetworks(ts.store, false)

	// then
	require.EqualError(t, err, "couldn't save network configuration: some error")
}

func testInitialisingConfigWithExistingConfigFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		NetworkExists("fairground").
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveNetwork(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := network.InitialiseNetworks(ts.store, false)

	// then
	require.Error(t, err)
}

func testOverwritingConfigSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		NetworkExists("fairground").
		Times(0)
	ts.store.EXPECT().
		SaveNetwork(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := network.InitialiseNetworks(ts.store, true)

	// then
	require.NoError(t, err)
}
