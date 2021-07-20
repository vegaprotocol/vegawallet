package v1_test

import (
	"os"
	"testing"

	crypto2 "code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Initialising new store succeeds", testInitialisingNewStoreSucceeds)
	t.Run("Saving already existing config fails", testFileStoreV1SavingAlreadyExistingConfigFails)
	t.Run("Saving new config  succeeds", testFileStoreV1SavingNewConfigSucceeds)
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

	_, err = os.Stat(configDir.RSAKeysPath())
	assert.NoError(t, err)
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

	// given
	cfg.Host = "my.new.host.com"

	// when
	err = s.SaveConfig(&cfg, true)

	// then
	require.NoError(t, err)

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

	// given
	newKeys := &service.RSAKeys{
		Pub:  []byte("my public key 2"),
		Priv: []byte("my private key 2"),
	}

	// when
	err = s.SaveRSAKeys(newKeys, true)

	// then
	require.NoError(t, err)

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

func newLegacyWalletWithKeys() *wallet.LegacyWallet {
	w := wallet.NewLegacyWallet("my-wallet")

	kp, err := wallet.GenKeyPair(crypto2.Ed25519, 1)
	if err != nil {
		panic(err)
	}

	w.KeyRing.Upsert(*kp)

	return w
}

func newHDWalletWithKeys() *wallet.HDWallet {
	w, _, err := wallet.NewHDWallet("my-wallet")
	if err != nil {
		panic(err)
	}

	_, err = w.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	return w
}
