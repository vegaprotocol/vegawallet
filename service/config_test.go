package service_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/mocks"
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
	t.Run("Generating config succeeds", testGeneratingConfigSucceeds)
	t.Run("Generating config with error fails", testGeneratingConfigWithErrorFails)
	t.Run("Generating config with RSA keys generation with error fails", testGeneratingConfigWithRSAKeysWithErrorFails)
	t.Run("Generating config with existing config fails", testGeneratingConfigWithExistingConfigFails)
	t.Run("Generating RSA keys with existing keys fails", testGeneratingRSAKeysWithExistingKeysFails)
	t.Run("Overwriting config and keys succeeds", testOverwritingConfigAndKeysSucceeds)
}

func testGeneratingConfigSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		RSAKeysExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := service.GenerateConfig(ts.store, false)

	// then
	require.NoError(t, err)
}

func testGeneratingConfigWithErrorFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(errors.New("some error"))
	ts.store.EXPECT().
		RSAKeysExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := service.GenerateConfig(ts.store, false)

	// then
	require.EqualError(t, err, "some error")
}

func testGeneratingConfigWithRSAKeysWithErrorFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		RSAKeysExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(1).
		Return(errors.New("some error"))

	// when
	err := service.GenerateConfig(ts.store, false)

	// then
	require.EqualError(t, err, "some error")
}

func testGeneratingConfigWithExistingConfigFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		RSAKeysExists().
		Times(0)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(0)

	// when
	err := service.GenerateConfig(ts.store, false)

	// then
	require.Error(t, err)
}

func testGeneratingRSAKeysWithExistingKeysFails(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(1).
		Return(false, nil)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		RSAKeysExists().
		Times(1).
		Return(true, nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(0)

	// when
	err := service.GenerateConfig(ts.store, false)

	// then
	require.Error(t, err)
}

func testOverwritingConfigAndKeysSucceeds(t *testing.T) {
	ts := getTestConfig(t)

	// setup
	ts.store.EXPECT().
		ConfigExists().
		Times(0)
	ts.store.EXPECT().
		SaveConfig(gomock.Any()).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		RSAKeysExists().
		Times(0)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any()).
		Times(1).
		Return(nil)

	// when
	err := service.GenerateConfig(ts.store, true)

	// then
	require.NoError(t, err)
}
