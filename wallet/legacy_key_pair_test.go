package wallet_test

import (
	"testing"

	crypto2 "code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyKeypair(t *testing.T) {
	t.Run("Generating a new key pair succeeds", testLegacyKeyPairGeneratingNewKeyPairSucceeds)
	t.Run("Generating a new key pair with unsupported algorithm fails", testLegacyKeyPairGeneratingNewKeyPairWithUnsupportedAlgorithmFails)
	t.Run("Tainting a key pair succeeds", testLegacyKeyPairTaintingKeyPairSucceeds)
	t.Run("Tainting an already tainted key pair fails", testLegacyKeyPairTaintingAlreadyTaintedKeyPairFails)
	t.Run("Secure copy of key pair removes sensitive information", testLegacyKeyPairToPublicKeyRemovesSensitiveInformation)
}

func testLegacyKeyPairGeneratingNewKeyPairSucceeds(t *testing.T) {
	// when
	kp, err := wallet.GenKeyPair(crypto2.Ed25519, 1)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, kp.Pub)
	assert.NotEmpty(t, kp.Priv)
	assert.False(t, kp.Tainted)
	assert.Equal(t, "vega/ed25519", kp.Algorithm.Name())
	assert.Equal(t, uint32(1), kp.Algorithm.Version())
	assert.Empty(t, kp.MetaList)
}

func testLegacyKeyPairGeneratingNewKeyPairWithUnsupportedAlgorithmFails(t *testing.T) {
	// when
	kp, err := wallet.GenKeyPair("unsupported-algo", 1)

	// then
	assert.Error(t, err)
	assert.Nil(t, kp)
}

func testLegacyKeyPairTaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.Tainted)
}

func testLegacyKeyPairTaintingAlreadyTaintedKeyPairFails(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.Tainted)

	// when
	err = kp.Taint()

	// then
	assert.Error(t, err)
	assert.True(t, kp.Tainted)
}

func testLegacyKeyPairToPublicKeyRemovesSensitiveInformation(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()

	// when
	secureKp := kp.ToPublicKey()

	// then
	assert.Equal(t, kp.Pub, secureKp.Pub)
	assert.Equal(t, kp.Tainted, secureKp.Tainted)
	assert.Equal(t, kp.Algorithm.Name(), secureKp.Algorithm.Name())
	assert.Equal(t, kp.Algorithm.Version(), secureKp.Algorithm.Version())
	assert.Equal(t, kp.MetaList, secureKp.MetaList)
}
