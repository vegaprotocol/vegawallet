package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeypair(t *testing.T) {
	t.Run("Generating a new key pair succeeds", testKeypairGeneratingNewKeyPairSucceeds)
	t.Run("Generating a new key pair with unsupported algorithm fails", testKeypairGeneratingNewKeyPairWithUnsupportedAlgorithmFails)
	t.Run("Tainting a key pair succeeds", testKeypairTaintingKeyPairSucceeds)
	t.Run("Tainting an already tainted key pair fails", testKeypairTaintingAlreadyTaintedKeyPairFails)
	t.Run("Secure copy of key pair removes sensitive information", testKeypairToPublicKeyRemovesSensitiveInformation)
}

func testKeypairGeneratingNewKeyPairSucceeds(t *testing.T) {
	// when
	kp, err := wallet.GenKeyPair(crypto.Ed25519)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, kp.Pub)
	assert.NotEmpty(t, kp.Priv)
	assert.False(t, kp.Tainted)
	assert.Equal(t, "vega/ed25519", kp.Algorithm.Name())
	assert.Equal(t, uint32(1), kp.Algorithm.Version())
	assert.Empty(t, kp.MetaList)
}

func testKeypairGeneratingNewKeyPairWithUnsupportedAlgorithmFails(t *testing.T) {
	// when
	kp, err := wallet.GenKeyPair("unsupported-algo")

	// then
	assert.Error(t, err)
	assert.Nil(t, kp)
}

func testKeypairTaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateKeyPair()

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.Tainted)
}

func testKeypairTaintingAlreadyTaintedKeyPairFails(t *testing.T) {
	// given
	kp := generateKeyPair()

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

func testKeypairToPublicKeyRemovesSensitiveInformation(t *testing.T) {
	// given
	kp := generateKeyPair()

	// when
	secureKp := kp.ToPublicKey()

	// then
	assert.Equal(t, kp.Pub, secureKp.Pub)
	assert.Equal(t, kp.Tainted, secureKp.Tainted)
	assert.Equal(t, kp.Algorithm.Name(), secureKp.Algorithm.Name())
	assert.Equal(t, kp.Algorithm.Version(), secureKp.Algorithm.Version())
	assert.Equal(t, kp.MetaList, secureKp.MetaList)
}
