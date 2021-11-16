package network_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/vegawallet/network"
	"code.vegaprotocol.io/vegawallet/network/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errSomethingWentWrong = errors.New("something went wrong")

type testConfig struct {
	ctrl  *gomock.Controller
	log   *zap.Logger
	store *mocks.MockStore
}

func getTestConfig(t *testing.T) *testConfig {
	t.Helper()
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	return &testConfig{
		ctrl:  ctrl,
		log:   zap.NewNop(),
		store: store,
	}
}

func TestImportNetwork(t *testing.T) {
	t.Run("Importing network succeeds", testImportingNetworkSucceeds)
	t.Run("Importing existing network fails", testImportingExistingNetworkFails)
	t.Run("Importing by overwriting existing network succeeds", testImportingByOverwritingNetworkSucceeds)
	t.Run("Importing network with errors when saving fails", testImportingNetworkWithErrorsWhenSavingFails)
}

func testImportingNetworkSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// given
	net := &network.Network{
		Name: "test",
	}

	// setup
	ts.store.EXPECT().
		NetworkExists("test").
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveNetwork(net).
		Times(1).
		Return(nil)

	// when
	err := network.ImportNetwork(ts.store, net, false)

	// then
	require.NoError(t, err)
}

func testImportingExistingNetworkFails(t *testing.T) {
	ts := getTestConfig(t)

	// given
	net := &network.Network{
		Name: "test",
	}

	// setup
	ts.store.EXPECT().
		NetworkExists("test").
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveNetwork(net).
		Times(0)

	// when
	err := network.ImportNetwork(ts.store, net, false)

	// then
	require.EqualError(t, err, "network \"test\" already exists")
}

func testImportingByOverwritingNetworkSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// given
	net := &network.Network{
		Name: "test",
	}

	// setup
	ts.store.EXPECT().
		NetworkExists("test").
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveNetwork(net).
		Times(1).
		Return(nil)

	// when
	err := network.ImportNetwork(ts.store, net, true)

	// then
	require.NoError(t, err)
}

func testImportingNetworkWithErrorsWhenSavingFails(t *testing.T) {
	ts := getTestConfig(t)

	// given
	net := &network.Network{
		Name: "test",
	}

	// setup
	ts.store.EXPECT().
		NetworkExists("test").
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveNetwork(net).
		Times(1).
		Return(errSomethingWentWrong)

	// when
	err := network.ImportNetwork(ts.store, net, true)

	// then
	require.EqualError(t, err, "couldn't save the imported network: something went wrong")
}
