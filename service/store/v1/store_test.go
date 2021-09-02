package v1_test

import (
	"io/fs"
	"os"
	"runtime"
	"testing"

	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/store/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Saving already existing config succeeds", testFileStoreV1SavingAlreadyExistingConfigSucceeds)
	t.Run("Saving new config succeeds", testFileStoreV1SavingNewConfigSucceeds)
	t.Run("Verifying non-existing config succeeds", testFileStoreV1VerifyingNonExistingConfigSucceeds)
	t.Run("Verifying config succeeds", testFileStoreV1VerifyingExistingConfigSucceeds)
	t.Run("Getting non-existing config fails", testFileStoreV1GetNonExistingConfigFails)
	t.Run("Getting config succeeds", testFileStoreV1GetConfigSucceeds)
	t.Run("Saving RSA keys already existing RSA keys succeeds", testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds)
	t.Run("Saving RSA keys succeeds", testFileStoreV1SaveRSAKeysSucceeds)
	t.Run("Verifying non-existing RSA keys fails", testFileStoreV1VerifyingNonExistingRSAKeysFails)
	t.Run("Verifying existing RSA keys succeeds", testFileStoreV1VerifyingExistingRSAKeysSucceeds)
	t.Run("Getting non-existing RSA keys fails", testFileStoreV1GetNonExistingRSAKeysFails)
	t.Run("Getting existing RSA keys succeeds", testFileStoreV1GetExistingRSAKeysSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	configDir := newVegaHome()

	s, err := v1.InitialiseStore(configDir.Paths())

	require.NoError(t, err)
	assert.NotNil(t, s)

	assertDirAccess(t, configDir.ConfigHome())
	assertDirAccess(t, configDir.RSAKeysHome())
}

func testFileStoreV1SavingAlreadyExistingConfigSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveConfig(&cfg)

	// then
	require.NoError(t, err)
}

func testFileStoreV1SavingNewConfigSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.ConfigFilePath())

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1VerifyingNonExistingConfigSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	exists, err := s.ConfigExists()

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingConfigSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg)

	// then
	require.NoError(t, err)

	// when
	exists, err := s.ConfigExists()

	// then
	assert.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingConfigFails(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	cfg, err := s.GetConfig()

	// then
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func testFileStoreV1GetConfigSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg)

	// then
	require.NoError(t, err)

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
}

func testFileStoreV1SaveRSAKeysSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.PublicRSAKeyFilePath())
	assertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func testFileStoreV1VerifyingNonExistingRSAKeysFails(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	exists, err := s.RSAKeysExists()

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.PublicRSAKeyFilePath())
	assertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// when
	exists, err := s.RSAKeysExists()

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingRSAKeysFails(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)

	// when
	keys, err := s.GetRsaKeys()

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newVegaHome()
	defer configDir.Remove()

	// given
	s := InitialiseFromPath(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.PublicRSAKeyFilePath())
	assertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func InitialiseFromPath(configDir vegaHome) *v1.Store {
	s, err := v1.InitialiseStore(configDir.Paths())
	if err != nil {
		panic(err)
	}

	return s
}

func assertDirAccess(t *testing.T, dirPath string) {
	stats, err := os.Stat(dirPath)
	require.NoError(t, err)
	assert.True(t, stats.IsDir())
	if runtime.GOOS == "windows" {
		assert.Equal(t, fs.FileMode(0777), stats.Mode().Perm())
	} else {
		assert.Equal(t, fs.FileMode(0700), stats.Mode().Perm())
	}
}

func assertFileAccess(t *testing.T, filePath string) {
	stats, err := os.Stat(filePath)
	assert.NoError(t, err)
	if runtime.GOOS == "windows" {
		assert.Equal(t, fs.FileMode(0666), stats.Mode().Perm())
	} else {
		assert.Equal(t, fs.FileMode(0600), stats.Mode().Perm())
	}
}
