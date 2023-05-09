package v1_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	vgtest "code.vegaprotocol.io/shared/libs/test"
	"code.vegaprotocol.io/vegawallet/wallet"
	storev1 "code.vegaprotocol.io/vegawallet/wallet/store/v1"

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
	walletsDir := newWalletsDir(t)

	s, err := storev1.InitialiseStore(walletsDir)

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, walletsDir)
}

func testFileStoreV1ListWalletsSucceeds(t *testing.T) {
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	passphrase := vgrand.RandomStr(5)

	var expectedWallets []string
	for i := 0; i < 3; i++ {
		w := newHDWalletWithKeys(t)

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
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	w := newHDWalletWithKeys(t)
	passphrase := vgrand.RandomStr(5)

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
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	w := newHDWalletWithKeys(t)
	passphrase := vgrand.RandomStr(5)
	othPassphrase := "not-original-passphrase"

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)

	// when
	returnedWallet, err := s.GetWallet(w.Name(), othPassphrase)

	// then
	assert.ErrorIs(t, err, wallet.ErrWrongPassphrase)
	assert.Nil(t, returnedWallet)
}

func testFileStoreV1GetNonExistingWalletFails(t *testing.T) {
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	name := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)

	// when
	returnedWallet, err := s.GetWallet(name, passphrase)

	// then
	assert.Error(t, err)
	assert.Nil(t, returnedWallet)
}

func testFileStoreV1GetWalletPathSucceeds(t *testing.T) {
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	name := vgrand.RandomStr(5)

	// when
	path := s.GetWalletPath(name)

	// then
	assert.Equal(t, filepath.Join(walletsDir, name), path)
}

func testFileStoreV1NonExistingWalletFails(t *testing.T) {
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	name := vgrand.RandomStr(5)

	// when
	exists := s.WalletExists(name)

	// then
	assert.False(t, exists)
}

func testFileStoreV1ExistingWalletSucceeds(t *testing.T) {
	walletsDir := newWalletsDir(t)

	// given
	s := initialiseStore(t, walletsDir)
	w := newHDWalletWithKeys(t)
	passphrase := vgrand.RandomStr(5)

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
	walletsDir := newWalletsDir(t)

	// given
	passphrase := vgrand.RandomStr(5)
	s := initialiseStore(t, walletsDir)
	w := newHDWalletWithKeys(t)

	// when
	err := s.SaveWallet(w, passphrase)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, filepath.Join(walletsDir, w.Name()))

	buf, err := ioutil.ReadFile(filepath.Join(walletsDir, w.Name()))
	if err != nil {
		t.Fatalf("couldn't read wallet file: %v", w.Name())
	}
	assert.NotEmpty(t, buf)
}

func initialiseStore(t *testing.T, walletsDir string) *storev1.Store {
	t.Helper()
	s, err := storev1.InitialiseStore(walletsDir)
	if err != nil {
		t.Fatalf("couldn't initialise store: %v", err)
	}

	return s
}

func newHDWalletWithKeys(t *testing.T) *wallet.HDWallet {
	t.Helper()
	w, _, err := wallet.NewHDWallet(fmt.Sprintf("my-wallet-%v-%s", time.Now().UnixNano(), vgrand.RandomStr(2)))
	if err != nil {
		t.Fatalf("couldn't create wallet: %v", err)
	}

	_, err = w.GenerateKeyPair([]wallet.Meta{})
	if err != nil {
		t.Fatalf("couldn't generate key: %v", err)
	}

	return w
}

func newWalletsDir(t *testing.T) string {
	t.Helper()
	rootPath := filepath.Join("/tmp", "vegawallet", vgrand.RandomStr(10))
	t.Cleanup(func() {
		if err := os.RemoveAll(rootPath); err != nil {
			t.Fatalf("couldn't remove vega home: %v", err)
		}
	})

	return rootPath
}
