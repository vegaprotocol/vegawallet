package v1_test

import (
	"testing"

	"code.vegaprotocol.io/vegawallet/service"
	"code.vegaprotocol.io/vegawallet/service/store/v1"
	vgtest "code.vegaprotocol.io/shared/libs/test"
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
	configDir := newVegaHome()
	defer configDir.Remove()

	s, err := v1.InitialiseStore(configDir.Paths())

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, configDir.RSAKeysHome())
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
	vgtest.AssertFileAccess(t, configDir.PublicRSAKeyFilePath())
	vgtest.AssertFileAccess(t, configDir.PrivateRSAKeyFilePath())

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
	vgtest.AssertFileAccess(t, configDir.PublicRSAKeyFilePath())
	vgtest.AssertFileAccess(t, configDir.PrivateRSAKeyFilePath())

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
	vgtest.AssertFileAccess(t, configDir.PublicRSAKeyFilePath())
	vgtest.AssertFileAccess(t, configDir.PrivateRSAKeyFilePath())

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func InitialiseFromPath(h vegaHome) *v1.Store {
	s, err := v1.InitialiseStore(h.Paths())
	if err != nil {
		panic(err)
	}

	return s
}
