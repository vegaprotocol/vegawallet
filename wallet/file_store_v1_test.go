package wallet_test

import (
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store without root folder fails", testNewFileStoreV1WithoutRootFolderFails)
	t.Run("New store with root folder succeeds and create wallet folder", testNewFileStoreV1WithRootFolderSucceeds)
	t.Run("Saving wallet succeeds", testFileStoreV1SaveWalletSucceeds)
	t.Run("Getting wallet succeeds", testFileStoreV1GetWalletSucceeds)
	t.Run("Getting wallet without wrong passphrase fails", testFileStoreV1GetWalletWithWrongPassphraseFails)
	t.Run("Getting non-existing wallet fails", testFileStoreV1GetNonExistingWalletFails)
	t.Run("Getting wallet path succeeds", testFileStoreV1GetWalletPathSucceeds)
}

func testNewFileStoreV1WithoutRootFolderFails(t *testing.T) {
	configDir := newConfigDir()

	s, err := wallet.NewFileStoreV1(configDir.RootPath())

	require.Error(t, err)
	assert.Nil(t, s)
}

func testNewFileStoreV1WithRootFolderSucceeds(t *testing.T) {
	configDir := newConfigDir()
	configDir.Create()
	defer configDir.Remove()

	s, err := wallet.NewFileStoreV1(configDir.RootPath())

	require.NoError(t, err)
	assert.NotNil(t, s)

	_, err = os.Stat(configDir.WalletsPath())
	assert.NoError(t, err)
}

func testFileStoreV1SaveWalletSucceeds(t *testing.T) {
	configDir := newConfigDir()
	configDir.Create()
	defer configDir.Remove()

	// given
	s := newFileStoreV1(configDir)
	w := newWalletWithKeys()

	// when
	err := s.SaveWallet(*w, "passphrase")

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, configDir.WalletContent(w.Owner))
}

func testFileStoreV1GetWalletSucceeds(t *testing.T) {
	configDir := newConfigDir()
	configDir.Create()
	defer configDir.Remove()

	// given
	s := newFileStoreV1(configDir)
	w := *newWalletWithKeys()
	passphrase := "passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Owner, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, w, returnedWallet)
}

func testFileStoreV1GetWalletWithWrongPassphraseFails(t *testing.T) {
	configDir := newConfigDir()
	configDir.Create()
	defer configDir.Remove()

	// given
	s := newFileStoreV1(configDir)
	w := *newWalletWithKeys()
	passphrase := "passphrase"
	othPassphrase := "not-original-passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Owner, othPassphrase)

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.Wallet{}, returnedWallet)
}

func testFileStoreV1GetNonExistingWalletFails(t *testing.T) {
	configDir := newConfigDir()
	configDir.Create()
	defer configDir.Remove()

	// given
	s := newFileStoreV1(configDir)
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
	configDir.Create()
	defer configDir.Remove()

	// given
	s := newFileStoreV1(configDir)
	name := "john"

	// when
	path := s.GetWalletPath(name)

	// then
	assert.Equal(t, configDir.WalletPath(name), path)
}

func newFileStoreV1(configDir configDir) *wallet.FileV1 {
	s, err := wallet.NewFileStoreV1(configDir.RootPath())
	if err != nil {
		panic(err)
	}
	return s
}

func newWalletWithKeys() *wallet.Wallet {
	w := wallet.NewWallet("my-wallet")

	kp, err := wallet.GenKeypair(crypto.Ed25519)
	if err != nil {
		panic(err)
	}

	w.KeyRing.Upsert(*kp)

	return w
}
