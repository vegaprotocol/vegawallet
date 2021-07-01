package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
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
	kp := generateKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// when
	ring.Upsert(*kp)

	// then
	assert.Contains(t, ring, *kp)
}

func testKeyRingUpdatingExistingKeySucceeds(t *testing.T) {
	// given
	kp := generateKeyPair()
	ring := wallet.NewLegacyKeyRing()

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
	kp1 := generateKeyPair()
	kp2 := generateKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	keys := ring.GetPublicKeys()

	// then
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, kp1.ToPublicKey())
	assert.Contains(t, keys, kp2.ToPublicKey())
}

func testKeyRingFindingExistingKeyPairSucceeds(t *testing.T) {
	// given
	kp1 := generateKeyPair()
	kp2 := generateKeyPair()
	ring := wallet.NewLegacyKeyRing()

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
	kp1 := generateKeyPair()
	kp2 := generateKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, err := ring.FindPair("non-existing public key")

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.LegacyKeyPair{}, returnedKp)
}
