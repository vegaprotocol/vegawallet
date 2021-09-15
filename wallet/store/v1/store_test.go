package v1_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"code.vegaprotocol.io/go-wallet/wallet"
	storev1 "code.vegaprotocol.io/go-wallet/wallet/store/v1"
	vgtest "code.vegaprotocol.io/shared/libs/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("Initialising store succeeds", testInitialisingStoreSucceeds)
	t.Run("Listing wallets succeeds", testFileStoreV1ListWalletsSucceeds)
	t.Run("Getting wallet succeeds", testFileStoreV1GetWalletSucceeds)
	t.Run("Getting wallet without wrong passphrase fails", testFileStoreV1GetWalletWithWrongPassphraseFails)
	t.Run("Getting non-existing wallet fails", testFileStoreV1GetNonExistingWalletFails)
	t.Run("Getting wallet path succeeds", testFileStoreV1GetWalletPathSucceeds)
	t.Run("Verifying non-existing wallet fails", testFileStoreV1NonExistingWalletFails)
	t.Run("Verifying existing wallet succeeds", testFileStoreV1ExistingWalletSucceeds)
	t.Run("Saving HD wallet succeeds", testFileStoreV1SaveHDWalletSucceeds)
}

func testInitialisingStoreSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()

	s, err := storev1.InitialiseStore(walletsDir.WalletsPath())

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, walletsDir.WalletsPath())
}

func testFileStoreV1ListWalletsSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := InitialiseStore(walletsDir)
	passphrase := "passphrase"

	var expectedWallets []string
	for i := 0; i < 3; i++ {
		w := newHDWalletWithKeys()

		// when
		err := s.SaveWallet(w, passphrase)

		// then
		require.NoError(t, err)

		expectedWallets = append(expectedWallets, w.Name())
	}
	sort.Strings(expectedWallets)

	// when
	returnedWallets, err := s.ListWallets()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedWallets, returnedWallets)
}

func testFileStoreV1GetWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	s := InitialiseStore(walletsDir)
	w := newHDWalletWithKeys()
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
	s := InitialiseStore(walletsDir)
	w := newHDWalletWithKeys()
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
	s := InitialiseStore(walletsDir)
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
	s := InitialiseStore(walletsDir)
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
	s := InitialiseStore(walletsDir)
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
	s := InitialiseStore(walletsDir)
	w := newHDWalletWithKeys()
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

func testFileStoreV1SaveHDWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir()
	defer walletsDir.Remove()

	// given
	passphrase := "passphrase"
	s := InitialiseStore(walletsDir)
	w := newHDWalletWithKeys()

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, walletsDir.WalletPath(w.Name()))
	assert.NotEmpty(t, walletsDir.WalletContent(w.Name()))
}

func InitialiseStore(walletsDir walletsDir) *storev1.Store {
	s, err := storev1.InitialiseStore(walletsDir.WalletsPath())
	if err != nil {
		panic(err)
	}

	return s
}

func newHDWalletWithKeys() *wallet.HDWallet {
	w, _, err := wallet.NewHDWallet(fmt.Sprintf("my-wallet-%v", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}

	_, err = w.GenerateKeyPair([]wallet.Meta{})
	if err != nil {
		panic(err)
	}

	return w
}
