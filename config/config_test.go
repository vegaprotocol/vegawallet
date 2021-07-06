package config_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/config"
	"code.vegaprotocol.io/go-wallet/config/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testService struct {
	ctrl  *gomock.Controller
	log   *zap.Logger
	store *mocks.MockStore
}

func getTestService(t *testing.T) *testService {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	return &testService{
		ctrl:  ctrl,
		log:   zap.NewNop(),
		store: store,
	}
}

func TestGenerateConfig(t *testing.T) {
	t.Run("Generating config succeeds", testGeneratingConfigSucceeds)
	t.Run("Generating config with error fails", testGeneratingConfigWithErrorFails)
	t.Run("Generating config with RSA keys generation succeeds", testGeneratingConfigWithRSAKeysSucceeds)
	t.Run("Generating config with RSA keys generation with error fails", testGeneratingConfigWithRSAKeysWithErrorFails)
	t.Run("Overwriting config succeeds", testOverwritingConfigSucceeds)
	t.Run("Overwriting config with RSA keys generation succeeds", testOverwritingConfigWithRSAKeysSucceeds)
}

func testGeneratingConfigSucceeds(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), false).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), gomock.Any()).
		Times(0)

	// when
	err := config.GenerateConfig(ts.log, ts.store, false, false)

	// then
	require.NoError(t, err)
}

func testGeneratingConfigWithErrorFails(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), false).
		Times(1).
		Return(errors.New("some error"))
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), gomock.Any()).
		Times(0)

	// when
	err := config.GenerateConfig(ts.log, ts.store, false, false)

	// then
	require.Error(t, err, errors.New("some error"))
}

func testGeneratingConfigWithRSAKeysSucceeds(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), false).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), false).
		Times(1).
		Return(nil)

	// when
	err := config.GenerateConfig(ts.log, ts.store, false, true)

	// then
	require.NoError(t, err)
}

func testGeneratingConfigWithRSAKeysWithErrorFails(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), false).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), false).
		Times(1).
		Return(errors.New("some error"))

	// when
	err := config.GenerateConfig(ts.log, ts.store, false, true)

	// then
	require.Error(t, err, errors.New("some error"))
}

func testOverwritingConfigSucceeds(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), true).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), gomock.Any()).
		Times(0)

	// when
	err := config.GenerateConfig(ts.log, ts.store, true, false)

	// then
	require.NoError(t, err)
}

func testOverwritingConfigWithRSAKeysSucceeds(t *testing.T) {
	ts := getTestService(t)

	// setup
	ts.store.EXPECT().
		SaveConfig(gomock.Any(), true).
		Times(1).
		Return(nil)
	ts.store.EXPECT().
		SaveRSAKeys(gomock.Any(), true).
		Times(1).
		Return(nil)

	// when
	err := config.GenerateConfig(ts.log, ts.store, true, true)

	// then
	require.NoError(t, err)
}
