package v1_test

import (
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/config"
	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Initialising new store succeeds", testInitialisingNewStoreSucceeds)
	t.Run("Saving wallet succeeds", testFileStoreV1SaveWalletSucceeds)
	t.Run("Getting wallet succeeds", testFileStoreV1GetWalletSucceeds)
	t.Run("Getting wallet without wrong passphrase fails", testFileStoreV1GetWalletWithWrongPassphraseFails)
	t.Run("Getting non-existing wallet fails", testFileStoreV1GetNonExistingWalletFails)
	t.Run("Getting wallet path succeeds", testFileStoreV1GetWalletPathSucceeds)
	t.Run("Verifying non-existing wallet fails", testFileStoreV1NonExistingWalletFails)
	t.Run("Verifying existing wallet succeeds", testFileStoreV1ExistingWalletSucceeds)
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

	s, err := storev1.NewStore(configDir.RootPath())

	require.NoError(t, err)
	assert.NotNil(t, s)
}

func testInitialisingNewStoreSucceeds(t *testing.T) {
	configDir := newConfigDir()

	s, err := storev1.NewStore(configDir.RootPath())

	require.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Initialise()

	_, err = os.Stat(configDir.WalletsPath())
	assert.NoError(t, err)

	_, err = os.Stat(configDir.RSAKeysPath())
	assert.NoError(t, err)
}

func testFileStoreV1SaveWalletSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	w := newWalletWithKeys()

	// when
	err := s.SaveWallet(*w, "passphrase")

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, configDir.WalletContent(w.Name))
}

func testFileStoreV1GetWalletSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	w := *newWalletWithKeys()
	passphrase := "passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, w, returnedWallet)
}

func testFileStoreV1GetWalletWithWrongPassphraseFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	w := *newWalletWithKeys()
	passphrase := "passphrase"
	othPassphrase := "not-original-passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Name, othPassphrase)

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.Wallet{}, returnedWallet)
}

func testFileStoreV1GetNonExistingWalletFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	name := "john"
	passphrase := "passphrase"

	// when
	returnedWallet, err := s.GetWallet(name, passphrase)

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.Wallet{}, returnedWallet)
}

func testFileStoreV1GetWalletPathSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	name := "john"

	// when
	path := s.GetWalletPath(name)

	// then
	assert.Equal(t, configDir.WalletPath(name), path)
}

func testFileStoreV1NonExistingWalletFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	name := "john"

	// when
	exists := s.WalletExists(name)

	// then
	assert.False(t, exists)
}

func testFileStoreV1ExistingWalletSucceeds(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	w := *newWalletWithKeys()
	passphrase := "passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	exists := s.WalletExists(w.Name)

	// then
	assert.True(t, exists)
}

func testFileStoreV1SavingAlreadyExistingConfigFails(t *testing.T) {
	configDir := newConfigDir()
	defer configDir.Remove()

	// given
	s := NewInitialisedStore(configDir)
	cfg := config.NewDefaultConfig()

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
	cfg := config.NewDefaultConfig()

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
	cfg := config.NewDefaultConfig()

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
	cfg := config.NewDefaultConfig()

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
	cfg := config.NewDefaultConfig()

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
	keys := &wallet.RSAKeys{
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
	keys := &wallet.RSAKeys{
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
	keys := &wallet.RSAKeys{
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
	keys := &wallet.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys, false)

	// then
	require.NoError(t, err)

	// given
	newKeys := &wallet.RSAKeys{
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
	keys := &wallet.RSAKeys{
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
	keys := &wallet.RSAKeys{
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

func NewStore(configDir configDir) *storev1.Store {
	s, err := storev1.NewStore(configDir.RootPath())
	if err != nil {
		panic(err)
	}

	return s
}

func NewInitialisedStore(configDir configDir) *storev1.Store {
	s := NewStore(configDir)

	err := s.Initialise()
	if err != nil {
		panic(err)
	}

	return s
}

func newWalletWithKeys() *wallet.Wallet {
	w := wallet.NewWallet("my-wallet")

	kp, err := wallet.GenKeyPair(crypto.Ed25519)
	if err != nil {
		panic(err)
	}

	w.KeyRing.Upsert(*kp)

	return w
}
