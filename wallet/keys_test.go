package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyRing(t *testing.T) {
	t.Run("Adding a new key succeeds", testKeyRingAddingNewKeySucceeds)
	t.Run("Updating an existing key succeeds", testKeyRingUpdatingExistingKeySucceeds)
	t.Run("Getting public keys succeeds", testKeyRingGettingPublicKeysSucceeds)
	t.Run("Finding an existing key pair succeeds", testKeyRingFindingExistingKeyPairSucceeds)
	t.Run("Finding a non-existing key pair fails", testKeyRingFindingNonExistingKeyPairFails)
}

func testKeyRingAddingNewKeySucceeds(t *testing.T) {
	// given
	kp := newKeyPair()
	ring := wallet.NewKeyRing()

	// when
	ring.Upsert(*kp)

	// then
	assert.Contains(t, ring, *kp)
}

func testKeyRingUpdatingExistingKeySucceeds(t *testing.T) {
	// given
	kp := newKeyPair()
	ring := wallet.NewKeyRing()

	// when
	ring.Upsert(*kp)

	// given
	updatedKp := *kp
	updatedKp.Tainted = true

	// when
	ring.Upsert(updatedKp)

	// then
	assert.Contains(t, ring, updatedKp)
	assert.NotContains(t, ring, *kp)
}

func testKeyRingGettingPublicKeysSucceeds(t *testing.T) {
	// given
	kp1 := newKeyPair()
	kp2 := newKeyPair()
	ring := wallet.NewKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	keys := ring.GetPublicKeys()

	// then
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, *kp1.ToPublicKey())
	assert.Contains(t, keys, *kp2.ToPublicKey())
}

func testKeyRingFindingExistingKeyPairSucceeds(t *testing.T) {
	// given
	kp1 := newKeyPair()
	kp2 := newKeyPair()
	ring := wallet.NewKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, err := ring.FindPair(kp2.Pub)

	// then
	require.NoError(t, err)
	assert.Equal(t, *kp2, returnedKp)
}

func testKeyRingFindingNonExistingKeyPairFails(t *testing.T) {
	// given
	kp1 := newKeyPair()
	kp2 := newKeyPair()
	ring := wallet.NewKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, err := ring.FindPair("non-existing public key")

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.KeyPair{}, returnedKp)
}

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
	assert.Empty(t, kp.Meta)
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
	kp := newKeyPair()

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.Tainted)
}

func testKeypairTaintingAlreadyTaintedKeyPairFails(t *testing.T) {
	// given
	kp := newKeyPair()

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
	kp := newKeyPair()

	// when
	secureKp := kp.ToPublicKey()

	// then
	assert.Equal(t, kp.Pub, secureKp.Key)
	assert.Equal(t, kp.Tainted, secureKp.Tainted)
	assert.Equal(t, kp.Algorithm.Name(), secureKp.Algorithm.Name())
	assert.Equal(t, kp.Algorithm.Version(), secureKp.Algorithm.Version())
	assert.Equal(t, kp.Meta, secureKp.Meta)
}

func newKeyPair() *wallet.KeyPair {
	kp, err := wallet.GenKeyPair(crypto.Ed25519)
	if err != nil {
		panic(err)
	}
	return kp
}
