package v1_test

import (
	"io/fs"
	"os"
	"runtime"
	"testing"

	"code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/wallet"
	storev1 "code.vegaprotocol.io/go-wallet/wallet/store/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Initialising new store succeeds", testInitialisingNewStoreSucceeds)
	t.Run("Getting wallet succeeds", testFileStoreV1GetWalletSucceeds)
	t.Run("Getting wallet without wrong passphrase fails", testFileStoreV1GetWalletWithWrongPassphraseFails)
	t.Run("Getting non-existing wallet fails", testFileStoreV1GetNonExistingWalletFails)
	t.Run("Getting wallet path succeeds", testFileStoreV1GetWalletPathSucceeds)
	t.Run("Verifying non-existing wallet fails", testFileStoreV1NonExistingWalletFails)
	t.Run("Verifying existing wallet succeeds", testFileStoreV1ExistingWalletSucceeds)
	t.Run("Saving legacy wallet succeeds", testFileStoreV1SaveLegacyWalletSucceeds)
	t.Run("Saving HD wallet succeeds", testFileStoreV1SaveHDWalletSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()

	s, err := storev1.NewStore(walletsDir.WalletsPath())

	require.NoError(t, err)
	assert.NotNil(t, s)
}

func testInitialisingNewStoreSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()

	s, err := storev1.NewStore(walletsDir.WalletsPath())

	require.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Initialise()

	assertDirAccess(t, walletsDir.WalletsPath())
}

func testFileStoreV1GetWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	w := newLegacyWalletWithKeys()
	passphrase := "passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Name(), passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, w, returnedWallet)
}

func testFileStoreV1GetWalletWithWrongPassphraseFails(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	w := newLegacyWalletWithKeys()
	passphrase := "passphrase"
	othPassphrase := "not-original-passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Name(), othPassphrase)

	// then
	assert.Error(t, err)
	assert.Equal(t, nil, returnedWallet)
}

func testFileStoreV1GetNonExistingWalletFails(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	name := "john"
	passphrase := "passphrase"

	// when
	returnedWallet, err := s.GetWallet(name, passphrase)

	// then
	assert.Error(t, err)
	assert.Equal(t, nil, returnedWallet)
}

func testFileStoreV1GetWalletPathSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	name := "john"

	// when
	path := s.GetWalletPath(name)

	// then
	assert.Equal(t, walletsDir.WalletPath(name), path)
}

func testFileStoreV1NonExistingWalletFails(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	name := "john"

	// when
	exists := s.WalletExists(name)

	// then
	assert.False(t, exists)
}

func testFileStoreV1ExistingWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	w := newLegacyWalletWithKeys()
	passphrase := "passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	exists := s.WalletExists(w.Name())

	// then
	assert.True(t, exists)
}

func testFileStoreV1SaveLegacyWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := NewInitialisedStore(walletsDir)
	w := newLegacyWalletWithKeys()

	// when
	err := s.SaveWallet(w, "passphrase")

	// then
	require.NoError(t, err)
	assertFileAccess(t, walletsDir.WalletPath(w.Name()))
	assert.NotEmpty(t, walletsDir.WalletContent(w.Name()))
}

func testFileStoreV1SaveHDWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	passphrase := "passphrase"
	s := NewInitialisedStore(walletsDir)
	w := newHDWalletWithKeys()

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)
	assertFileAccess(t, walletsDir.WalletPath(w.Name()))
	assert.NotEmpty(t, walletsDir.WalletContent(w.Name()))
}

func NewStore(walletsDir walletsDir) *storev1.Store {
	s, err := storev1.NewStore(walletsDir.WalletsPath())
	if err != nil {
		panic(err)
	}

	return s
}

func NewInitialisedStore(walletsDir walletsDir) *storev1.Store {
	s := NewStore(walletsDir)

	err := s.Initialise()
	if err != nil {
		panic(err)
	}

	return s
}

func newLegacyWalletWithKeys() *wallet.LegacyWallet {
	w := wallet.NewLegacyWallet("my-wallet")

	kp, err := wallet.GenKeyPair(crypto.Ed25519, 1)
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
