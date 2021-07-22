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
	t.Run("Initialising new store succeeds", testInitialisingNewStoreSucceeds)
	t.Run("Saving already existing config fails", testFileStoreV1SavingAlreadyExistingConfigFails)
	t.Run("Saving new config succeeds", testFileStoreV1SavingNewConfigSucceeds)
	t.Run("Overwriting existing config succeeds", testFileStoreV1OverwritingExistingConfigSucceeds)
	t.Run("Overwriting non-existing config succeeds", testFileStoreV1OverwritingNonExistingConfigSucceeds)
	t.Run("Getting non-existing config fails", testFileStoreV1GetNonExistingConfigFails)
	t.Run("Getting config succeeds", testFileStoreV1GetConfigSucceeds)
	t.Run("Saving RSA keys without folder fails", testFileStoreV1SaveRSAKeysWithoutFolderFails)
	t.Run("Saving RSA keys already existing RSA keys fails", testFileStoreV1SaveAlreadyExistingRSAKeysFails)
	t.Run("Saving RSA keys succeeds", testFileStoreV1SaveRSAKeysSucceeds)
	t.Run("Overwriting already existing RSA keys succeeds", testFileStoreV1OverwritingAlreadyExistingRSAKeysSucceeds)
	t.Run("Overwriting non-existing RSA keys succeeds", testFileStoreV1OverwritingNonExistingRSAKeysSucceeds)
	t.Run("Getting non-existing RSA keys fails", testFileStoreV1GetNonExistingRSAKeysFails)
	t.Run("Getting existing RSA keys succeeds", testFileStoreV1GetExistingRSAKeysSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	configDir := newConfigDir()

	s, err := v1.NewStore(configDir.RootPath())

	require.NoError(t, err)
	assert.NotNil(t, s)
}

func testInitialisingNewStoreSucceeds(t *testing.T) {
	configDir := newConfigDir()

	s, err := v1.NewStore(configDir.RootPath())

	require.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Initialise()

	assertDirAccess(t, configDir.RootPath())
	assertDirAccess(t, configDir.RSAKeysPath())
}

func testFileStoreV1SavingAlreadyExistingConfigFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg, false)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveConfig(&cfg, false)

	// then
	require.Error(t, err)
}

func testFileStoreV1SavingNewConfigSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg, false)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.ConfigFilePath())

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1OverwritingNonExistingConfigSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg, true)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.ConfigFilePath())

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1OverwritingExistingConfigSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg, false)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.ConfigFilePath())

	// given
	cfg.Host = "my.new.host.com"

	// when
	err = s.SaveConfig(&cfg, true)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.ConfigFilePath())

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1GetNonExistingConfigFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)

	// when
	cfg, err := s.GetConfig()

	// then
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func testFileStoreV1GetConfigSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := service.NewDefaultConfig()

	// when
	err := s.SaveConfig(&cfg, false)

	// then
	require.NoError(t, err)

	// when
	returnedCfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, &cfg, returnedCfg)
}

func testFileStoreV1SaveRSAKeysWithoutFolderFails(t *testing.T) {
	configDir := newConfigDir()

	// given
	s := NewStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

	// then
	require.Error(t, err)
}

func testFileStoreV1SaveAlreadyExistingRSAKeysFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveRSAKeys(keys, false)

	// then
	require.Error(t, err)
}

func testFileStoreV1SaveRSAKeysSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

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

func testFileStoreV1OverwritingAlreadyExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

	// then
	require.NoError(t, err)
	assertFileAccess(t, configDir.PublicRSAKeyFilePath())
	assertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// given
	newKeys := &service.RSAKeys{
		Pub:  []byte("my public key 2"),
		Priv: []byte("my private key 2"),
	}

	// when
	err = s.SaveRSAKeys(newKeys, true)

	// then
	assertFileAccess(t, configDir.PublicRSAKeyFilePath())
	assertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, newKeys, returnedKeys)
}

func testFileStoreV1OverwritingNonExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, true)

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

func testFileStoreV1GetNonExistingRSAKeysFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)

	// when
	keys, err := s.GetRsaKeys()

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingRSAKeysSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

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

func NewStore(configDir configDir) *v1.Store {
	s, err := v1.NewStore(configDir.RootPath())
	if err != nil {
		panic(err)
	}

	return s
}

func NewInitialisedStore(configDir configDir) *v1.Store {
	s := NewStore(configDir)

	err := s.Initialise()
	if err != nil {
		panic(err)
	}

	return s
}

func assertDirAccess(t *testing.T, dirPath string) {
	stats, err := os.Stat(dirPath)
	assert.NoError(t, err)
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
