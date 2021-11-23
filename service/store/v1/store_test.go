package v1_test

import (
	"os"
	"path/filepath"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	vgtest "code.vegaprotocol.io/shared/libs/test"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/service"
	v1 "code.vegaprotocol.io/vegawallet/service/store/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Saving already existing RSA keys succeeds", testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds)
	t.Run("Saving RSA keys succeeds", testFileStoreV1SaveRSAKeysSucceeds)
	t.Run("Verifying non-existing RSA keys fails", testFileStoreV1VerifyingNonExistingRSAKeysFails)
	t.Run("Verifying existing RSA keys succeeds", testFileStoreV1VerifyingExistingRSAKeysSucceeds)
	t.Run("Getting non-existing RSA keys fails", testFileStoreV1GetNonExistingRSAKeysFails)
	t.Run("Getting existing RSA keys succeeds", testFileStoreV1GetExistingRSAKeysSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	s, err := v1.InitialiseStore(vegaHome)

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, rsaKeysHome(t, vegaHome))
}

func testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
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
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, vegaHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, vegaHome))

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func testFileStoreV1VerifyingNonExistingRSAKeysFails(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	exists, err := s.RSAKeysExists()

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingRSAKeysSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, vegaHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, vegaHome))

	// when
	exists, err := s.RSAKeysExists()

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingRSAKeysFails(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)

	// when
	keys, err := s.GetRsaKeys()

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingRSAKeysSucceeds(t *testing.T) {
	vegaHome := newVegaHome(t)

	// given
	s := initialiseFromPath(t, vegaHome)
	keys := &service.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, vegaHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, vegaHome))

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func initialiseFromPath(t *testing.T, vegaHome *paths.CustomPaths) *v1.Store {
	t.Helper()
	s, err := v1.InitialiseStore(vegaHome)
	if err != nil {
		t.Fatalf("couldn't initialise store: %v", err)
	}

	return s
}

func newVegaHome(t *testing.T) *paths.CustomPaths {
	t.Helper()
	rootPath := filepath.Join("/tmp", "vegawallet", vgrand.RandomStr(10))
	t.Cleanup(func() {
		if err := os.RemoveAll(rootPath); err != nil {
			t.Fatalf("couldn't remove vega home: %v", err)
		}
	})

	return &paths.CustomPaths{CustomHome: rootPath}
}

func rsaKeysHome(t *testing.T, vegaHome *paths.CustomPaths) string {
	t.Helper()
	return vegaHome.DataPathFor(paths.WalletServiceRSAKeysDataHome)
}

func publicRSAKeyFilePath(t *testing.T, vegaHome *paths.CustomPaths) string {
	t.Helper()
	return vegaHome.DataPathFor(paths.WalletServicePublicRSAKeyDataFile)
}

func privateRSAKeyFilePath(t *testing.T, vegaHome *paths.CustomPaths) string {
	t.Helper()
	return vegaHome.DataPathFor(paths.WalletServicePrivateRSAKeyDataFile)
}
